package services

import (
	"context"
	"errors"
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
	return s.run(ctx, ownerID, prompt, false, nil)
}

func (s *AIBuilderService) Refine(ctx context.Context, ownerID, prompt string, current *models.StorefrontConfig) (DesignResult, error) {
	return s.run(ctx, ownerID, prompt, true, current)
}

func (s *AIBuilderService) run(ctx context.Context, ownerID, prompt string, refine bool, current *models.StorefrontConfig) (DesignResult, error) {
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
	base := storefront.Normalize(store.Draft)
	if refine && current != nil {
		base = storefront.Normalize(*current)
	}
	prior := priorTurns(store.DesignNote)
	storeCtx.History = historyFor(prior)
	var result ai.Result
	mode := "generate"
	if refine {
		mode = "refine"
		result, err = s.provider.Refine(ctx, prompt, base, storeCtx)
	} else {
		result, err = s.provider.Generate(ctx, prompt, storeCtx)
	}
	if err != nil {
		message := "The designer could not update the storefront. Try again."
		if detail := ai.SafeDetail(err); detail != "" {
			message = detail
		}
		return DesignResult{}, apperr.New(http.StatusBadGateway, "design_failed", message)
	}
	cfg := storefront.Normalize(result.Config)
	summary := result.Summary
	if result.Provider == "mock" && !strings.Contains(strings.ToLower(summary), "mock") {
		summary = "Mock designer: " + summary
	}
	note := &models.DesignNote{
		Provider: result.Provider, Mode: mode, Prompt: prompt, Summary: summary,
		Turns: appendTurns(prior, prompt, summary),
	}
	now := time.Now().UTC()
	if err := s.stores.UpdateDraft(ctx, ownerID, store.ID, cfg, note, now); err != nil {
		return DesignResult{}, err
	}
	store.Draft = cfg
	store.DesignNote = note
	store.UpdatedAt = now
	if store.Published {
		if err := s.stores.PublishDraft(ctx, ownerID, store.ID, now); err != nil {
			if errors.Is(err, interfaces.ErrNotFound) {
				return DesignResult{}, apperr.New(http.StatusNotFound, "not_found", "Store not found")
			}
			return DesignResult{}, err
		}
		store.Storefront = cfg
	}
	return DesignResult{Store: store, Provider: result.Provider, Summary: summary}, nil
}

func priorTurns(note *models.DesignNote) []models.DesignTurn {
	if note == nil {
		return nil
	}
	if len(note.Turns) > 0 {
		return note.Turns
	}
	var turns []models.DesignTurn
	if text := clipTurn(note.Prompt); text != "" {
		turns = append(turns, models.DesignTurn{Role: "user", Text: text})
	}
	if text := clipTurn(note.Summary); text != "" {
		turns = append(turns, models.DesignTurn{Role: "designer", Text: text})
	}
	return turns
}

func appendTurns(prior []models.DesignTurn, prompt, summary string) []models.DesignTurn {
	turns := append([]models.DesignTurn{}, prior...)
	turns = append(turns, models.DesignTurn{Role: "user", Text: clipTurn(prompt)}, models.DesignTurn{Role: "designer", Text: clipTurn(summary)})
	if len(turns) > 16 {
		turns = turns[len(turns)-16:]
	}
	return turns
}

func historyFor(turns []models.DesignTurn) []ai.Turn {
	if len(turns) == 0 {
		return nil
	}
	if len(turns) > 12 {
		turns = turns[len(turns)-12:]
	}
	out := make([]ai.Turn, 0, len(turns))
	for _, turn := range turns {
		out = append(out, ai.Turn{Role: turn.Role, Text: turn.Text})
	}
	return out
}

func clipTurn(value string) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) > 400 {
		return string(runes[:400])
	}
	return value
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
