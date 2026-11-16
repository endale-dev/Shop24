package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MongoURI     string
	MongoDB      string
	JWTSecret    string
	JWTTTL       time.Duration
	CookieSecure bool
	CookieDomain string
	CORSOrigins  []string
	BaseDomain   string
	AIProvider   string
	AIAPIKey     string
	AIModel      string
	AIBaseURL    string
	AITimeout    time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()

	ttlHours, err := strconv.Atoi(env("JWT_TTL_HOURS", "168"))
	if err != nil || ttlHours < 1 || ttlHours > 24*30 {
		return Config{}, fmt.Errorf("JWT_TTL_HOURS must be between 1 and 720")
	}
	aiTimeout, err := strconv.Atoi(env("AI_TIMEOUT_SECONDS", "30"))
	if err != nil || aiTimeout < 5 || aiTimeout > 120 {
		return Config{}, fmt.Errorf("AI_TIMEOUT_SECONDS must be between 5 and 120")
	}

	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 16 {
		return Config{}, fmt.Errorf("JWT_SECRET must be set and at least 16 characters")
	}

	provider := strings.ToLower(env("AI_PROVIDER", "mock"))
	if provider != "mock" && provider != "openai" {
		return Config{}, fmt.Errorf("AI_PROVIDER must be mock or openai")
	}

	origins := splitCSV(env("CORS_ORIGINS", "http://localhost:5173"))
	base := strings.ToLower(strings.TrimSpace(env("BASE_DOMAIN", "localhost")))
	if base == "" {
		return Config{}, fmt.Errorf("BASE_DOMAIN is required")
	}

	return Config{
		Port:         env("PORT", "8080"),
		MongoURI:     env("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      env("MONGO_DB", "sable"),
		JWTSecret:    secret,
		JWTTTL:       time.Duration(ttlHours) * time.Hour,
		CookieSecure: env("COOKIE_SECURE", "false") == "true",
		CookieDomain: strings.TrimSpace(os.Getenv("COOKIE_DOMAIN")),
		CORSOrigins:  origins,
		BaseDomain:   base,
		AIProvider:   provider,
		AIAPIKey:     strings.TrimSpace(os.Getenv("AI_API_KEY")),
		AIModel:      env("AI_MODEL", "gpt-4o-mini"),
		AIBaseURL:    strings.TrimRight(env("AI_BASE_URL", "https://api.openai.com/v1"), "/"),
		AITimeout:    time.Duration(aiTimeout) * time.Second,
	}, nil
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
