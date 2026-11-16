package middleware

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/httpx"
	"shop24/internal/token"
)

func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic", "error", rec, "path", r.URL.Path)
					httpx.Error(w, apperr.New(http.StatusInternalServerError, "internal", "Something went wrong"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = token.NewID()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func Log(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			writer := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(writer, r)
			log.Info("request", "method", r.Method, "path", r.URL.Path, "status", writer.status, "ms", time.Since(start).Milliseconds())
		})
	}
}

func CORS(origins []string, baseDomain string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, origin := range origins {
		allowed[origin] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && originAllowed(origin, allowed, baseDomain) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-Id")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				if origin != "" && !originAllowed(origin, allowed, baseDomain) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(issuer *token.Issuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(token.CookieName)
			if err != nil || cookie.Value == "" {
				httpx.Error(w, apperr.ErrUnauthorized)
				return
			}
			userID, err := issuer.Parse(cookie.Value)
			if err != nil {
				httpx.Error(w, apperr.ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(token.WithUserID(r.Context(), userID)))
		})
	}
}

func originAllowed(origin string, allowed map[string]struct{}, baseDomain string) bool {
	if _, ok := allowed[origin]; ok {
		return true
	}
	parsed, err := url.Parse(origin)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	base := strings.ToLower(baseDomain)
	return host == "localhost" || host == "127.0.0.1" || host == base || strings.HasSuffix(host, "."+base)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
