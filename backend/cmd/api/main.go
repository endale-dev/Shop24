package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shop24/internal/ai"
	"shop24/internal/config"
	"shop24/internal/handlers"
	mongorepo "shop24/internal/repositories/mongo"
	"shop24/internal/server"
	"shop24/internal/services"
	"shop24/internal/token"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg, err := config.Load()
	if err != nil {
		log.Error("configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := mongorepo.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Error("mongodb", "error", err)
		os.Exit(1)
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(disconnectCtx); err != nil {
			log.Error("mongodb disconnect", "error", err)
		}
	}()

	db := client.Database(cfg.MongoDB)
	if err := mongorepo.EnsureIndexes(ctx, db); err != nil {
		log.Error("mongodb indexes", "error", err)
		os.Exit(1)
	}

	users := mongorepo.NewUserRepository(db)
	stores := mongorepo.NewStoreRepository(db)
	products := mongorepo.NewProductRepository(db)
	collections := mongorepo.NewCollectionRepository(db)
	orders := mongorepo.NewOrderRepository(db)
	issuer := token.NewIssuer(cfg.JWTSecret, cfg.JWTTTL)

	authService := services.NewAuthService(users, issuer)
	storeService := services.NewStoreService(stores)
	productService := services.NewProductService(stores, products, collections)
	collectionService := services.NewCollectionService(stores, products, collections)
	orderService := services.NewOrderService(stores, products, orders)
	analyticsService := services.NewAnalyticsService(stores, products, orders)
	storefrontService := services.NewStorefrontService(stores)

	var provider ai.Provider
	switch cfg.AIProvider {
	case "openai":
		if cfg.AIAPIKey == "" {
			log.Warn("AI_PROVIDER is openai but AI_API_KEY is empty; using the mock designer")
			provider = ai.NewMock()
		} else {
			provider = ai.NewOpenAI(cfg.AIAPIKey, cfg.AIBaseURL, cfg.AIModel, cfg.AITimeout)
		}
	default:
		provider = ai.NewMock()
	}
	log.Info("storefront designer", "provider", provider.Name())
	builder := services.NewAIBuilderService(stores, products, provider)

	srv := server.New(server.App{
		Config:      cfg,
		Log:         log,
		Tokens:      issuer,
		Auth:        handlers.NewAuthHandler(authService, log, cfg.CookieSecure, cfg.CookieDomain),
		Stores:      handlers.NewStoreHandler(storeService, log),
		Products:    handlers.NewProductHandler(productService, log),
		Collections: handlers.NewCollectionHandler(collectionService, log),
		Orders:      handlers.NewOrderHandler(orderService, log),
		Analytics:   handlers.NewAnalyticsHandler(analyticsService, log),
		Storefront:  handlers.NewStorefrontHandler(storefrontService, builder, log),
		Public:      handlers.NewPublicHandler(storeService, productService, orderService, collectionService, log),
		Health:      func(ctx context.Context) error { return client.Ping(ctx, nil) },
	})
	if err := srv.Run(ctx); err != nil {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
