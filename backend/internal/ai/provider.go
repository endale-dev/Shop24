package ai

import (
	"context"
	"strings"

	"shop24/internal/domain/models"
)

type ProviderError struct {
	Detail string
}

func (e *ProviderError) Error() string {
	return e.Detail
}

func SafeDetail(err error) string {
	if err == nil {
		return ""
	}
	detail := strings.TrimSpace(err.Error())
	if detail == "" || strings.Contains(strings.ToLower(detail), "sk-") {
		return ""
	}
	runes := []rune(detail)
	if len(runes) > 160 {
		detail = string(runes[:160])
	}
	return detail
}

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
