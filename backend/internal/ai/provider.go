package ai

import (
	"context"

	"shop24/internal/domain/models"
)

type StoreContext struct {
	Name         string
	Description  string
	Categories   []string
	ProductCount int
}

type Result struct {
	Config   models.StorefrontConfig
	Provider string
	Summary  string
}

// Provider turns a natural-language prompt into a storefront configuration.
// Implementations must not return executable code.
type Provider interface {
	Name() string
	Generate(ctx context.Context, prompt string, store StoreContext) (Result, error)
	Refine(ctx context.Context, prompt string, current models.StorefrontConfig, store StoreContext) (Result, error)
}
