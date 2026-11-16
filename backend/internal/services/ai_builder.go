package services

import (
	"context"
	"net/http"
	"strings"
	"time"

	"shop24/internal/ai"
	"shop24/internal/apperr"
	"shop24/internal/domain/models"
	"shop24/internal/domain/storefront"
	"shop24/internal/repositories/interfaces"
)

type AIBuilderService struct {
	stores   interfaces.StoreRepository
	products interfaces.ProductRepository
	provider ai.Provider
}

func NewAIBuilderService(stores interfaces.StoreRepository, products interfaces.ProductRepository, provider ai.Provider) *AIBuilderService {
	return &AIBuilderService{stores: stores, products: products, provider: provider}
}

type DesignResult struct {
	Store    models.Store `json:"store"`
	Provider string       `json:"provider"`
	Summary  string       `json:"summary"`
}

func (s *AIBuilderService) Generate(ctx context.Context, ownerID, prompt string) (DesignResult, error) {
	return s.run(ctx, ownerID, prompt, false)
}

func (s *AIBuilderService) Refine(ctx context.Context, ownerID, prompt string) (DesignResult, error) {
	return s.run(ctx, ownerID, prompt, true)
}

func (s *AIBuilderService) run(ctx context.Context, ownerID, prompt string, refine bool) (DesignResult, error) {
	prompt = strings.TrimSpace(prompt)
	if len([]rune(prompt)) < 8 || len([]rune(prompt)) > 1000 {
		return DesignResult{}, apperr.New(http.StatusBadRequest, "invalid_prompt", "Describe the design in 8 to 1000 characters")
	}
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return DesignResult{}, err
	}
	storeCtx, err := s.contextFor(ctx, store)
	if err != nil {
		return DesignResult{}, err
	}
	var result ai.Result
	mode := "generate"
	if refine {
		mode = "refine"
		result, err = s.provider.Refine(ctx, prompt, storefront.Normalize(store.Draft), storeCtx)
	} else {
		result, err = s.provider.Generate(ctx, prompt, storeCtx)
	}
	if err != nil {
		return DesignResult{}, apperr.New(http.StatusBadGateway, "design_failed", "The designer could not update the storefront. Try again.")
	}
	cfg := storefront.Normalize(result.Config)
	summary := result.Summary
	if result.Provider == "mock" && !strings.Contains(strings.ToLower(summary), "mock") {
		summary = "Mock designer: " + summary
	}
	note := &models.DesignNote{Provider: result.Provider, Mode: mode, Prompt: prompt, Summary: summary}
	now := time.Now().UTC()
	if err := s.stores.UpdateDraft(ctx, ownerID, store.ID, cfg, note, now); err != nil {
		return DesignResult{}, err
	}
	store.Draft = cfg
	store.DesignNote = note
	store.UpdatedAt = now
	return DesignResult{Store: store, Provider: result.Provider, Summary: summary}, nil
}

func (s *AIBuilderService) contextFor(ctx context.Context, store models.Store) (ai.StoreContext, error) {
	products, err := s.products.List(ctx, store.ID, interfaces.ProductFilter{Limit: 200})
	if err != nil {
		return ai.StoreContext{}, err
	}
	seen := map[string]struct{}{}
	categories := make([]string, 0)
	for _, product := range products {
		if product.Category == "" {
			continue
		}
		if _, ok := seen[product.Category]; ok {
			continue
		}
		seen[product.Category] = struct{}{}
		categories = append(categories, product.Category)
		if len(categories) == 12 {
			break
		}
	}
	return ai.StoreContext{
		Name: store.Name, Description: store.Description, Categories: categories, ProductCount: len(products),
	}, nil
}
