package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/config"
	"shop24/internal/handlers"
	"shop24/internal/httpx"
	"shop24/internal/middleware"
	"shop24/internal/token"
)

type App struct {
	Config      config.Config
	Log         *slog.Logger
	Tokens      *token.Issuer
	Auth        *handlers.AuthHandler
	Stores      *handlers.StoreHandler
	Products    *handlers.ProductHandler
	Collections *handlers.CollectionHandler
	Orders      *handlers.OrderHandler
	Analytics   *handlers.AnalyticsHandler
	Storefront  *handlers.StorefrontHandler
	Public      *handlers.PublicHandler
	Health      func(context.Context) error
}

type Server struct {
	http *http.Server
	log  *slog.Logger
}

func New(app App) *Server {
	mux := http.NewServeMux()
	auth := middleware.RequireAuth(app.Tokens)

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := app.Health(ctx); err != nil {
			app.Log.Error("health check failed", "error", err)
			httpx.Error(w, apperr.New(http.StatusServiceUnavailable, "unavailable", "Database unavailable"))
			return
		}
		httpx.Data(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /api/auth/register", app.Auth.Register)
	mux.HandleFunc("POST /api/auth/login", app.Auth.Login)
	mux.HandleFunc("POST /api/auth/logout", app.Auth.Logout)
	mux.HandleFunc("GET /api/auth/providers", app.Auth.Providers)
	mux.HandleFunc("GET /api/auth/google/start", app.Auth.GoogleStart)
	mux.HandleFunc("GET /api/auth/google/callback", app.Auth.GoogleCallback)
	mux.HandleFunc("POST /api/auth/forgot", app.Auth.Forgot)
	mux.HandleFunc("POST /api/auth/reset", app.Auth.Reset)
	mux.Handle("GET /api/auth/me", auth(http.HandlerFunc(app.Auth.Me)))
	mux.Handle("PATCH /api/auth/profile", auth(http.HandlerFunc(app.Auth.UpdateProfile)))
	mux.Handle("PATCH /api/auth/password", auth(http.HandlerFunc(app.Auth.ChangePassword)))

	mux.Handle("POST /api/stores", auth(http.HandlerFunc(app.Stores.Create)))
	mux.Handle("GET /api/stores/mine", auth(http.HandlerFunc(app.Stores.Mine)))
	mux.Handle("PATCH /api/stores/mine", auth(http.HandlerFunc(app.Stores.Update)))
	mux.Handle("POST /api/stores/mine/publish", auth(http.HandlerFunc(app.Stores.Publish)))
	mux.Handle("POST /api/stores/mine/unpublish", auth(http.HandlerFunc(app.Stores.Unpublish)))

	mux.Handle("GET /api/products", auth(http.HandlerFunc(app.Products.List)))
	mux.Handle("POST /api/products", auth(http.HandlerFunc(app.Products.Create)))
	mux.Handle("GET /api/products/{id}", auth(http.HandlerFunc(app.Products.Get)))
	mux.Handle("PATCH /api/products/{id}", auth(http.HandlerFunc(app.Products.Update)))
	mux.Handle("DELETE /api/products/{id}", auth(http.HandlerFunc(app.Products.Delete)))
	mux.Handle("POST /api/products/{id}/stock", auth(http.HandlerFunc(app.Products.SetStock)))

	mux.Handle("GET /api/collections", auth(http.HandlerFunc(app.Collections.List)))
	mux.Handle("POST /api/collections", auth(http.HandlerFunc(app.Collections.Create)))
	mux.Handle("GET /api/collections/{id}", auth(http.HandlerFunc(app.Collections.Get)))
	mux.Handle("PATCH /api/collections/{id}", auth(http.HandlerFunc(app.Collections.Update)))
	mux.Handle("DELETE /api/collections/{id}", auth(http.HandlerFunc(app.Collections.Delete)))

	mux.Handle("GET /api/orders", auth(http.HandlerFunc(app.Orders.List)))
	mux.Handle("GET /api/orders/{id}", auth(http.HandlerFunc(app.Orders.Get)))
	mux.Handle("PATCH /api/orders/{id}", auth(http.HandlerFunc(app.Orders.UpdateStatus)))

	mux.Handle("GET /api/analytics", auth(http.HandlerFunc(app.Analytics.Overview)))

	mux.Handle("POST /api/storefront/generate", auth(http.HandlerFunc(app.Storefront.Generate)))
	mux.Handle("POST /api/storefront/refine", auth(http.HandlerFunc(app.Storefront.Refine)))
	mux.Handle("PUT /api/storefront", auth(http.HandlerFunc(app.Storefront.Save)))
	mux.Handle("POST /api/storefront/publish", auth(http.HandlerFunc(app.Storefront.Publish)))

	mux.HandleFunc("GET /api/public/stores/{slug}", app.Public.Store)
	mux.HandleFunc("GET /api/public/stores/{slug}/products", app.Public.Products)
	mux.HandleFunc("GET /api/public/stores/{slug}/products/{id}", app.Public.Product)
	mux.HandleFunc("GET /api/public/stores/{slug}/collections", app.Public.Collections)
	mux.HandleFunc("GET /api/public/stores/{slug}/collections/{collectionSlug}", app.Public.Collection)
	mux.HandleFunc("POST /api/public/stores/{slug}/orders", app.Public.PlaceOrder)

	handler := middleware.Recover(app.Log)(middleware.RequestID(middleware.Log(app.Log)(middleware.CORS(app.Config.CORSOrigins, app.Config.BaseDomain)(mux))))
	return &Server{
		log: app.Log,
		http: &http.Server{
			Addr:              ":" + app.Config.Port,
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      90 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		s.log.Info("api listening", "addr", s.http.Addr)
		errCh <- s.http.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.http.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
