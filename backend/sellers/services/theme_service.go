package services

import (
	"context"
	"time"

	"github.com/Endale2/DRPS/sellers/models"
	"github.com/Endale2/DRPS/sellers/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThemeService struct {
	repo *repositories.ThemeRepository
}

func NewThemeService() *ThemeService {
	return &ThemeService{repo: repositories.NewThemeRepository()}
}

func (s *ThemeService) Create(ctx context.Context, t *models.Theme) error {
	return s.repo.Create(ctx, t)
}

func (s *ThemeService) ListByStore(ctx context.Context, storeID primitive.ObjectID) ([]models.Theme, error) {
	return s.repo.FindByStore(ctx, storeID)
}

func (s *ThemeService) Activate(ctx context.Context, storeID primitive.ObjectID, themeID primitive.ObjectID) error {
	return s.repo.SetActive(ctx, storeID, themeID)
}

// ThemeModelFromMap builds a models.Theme from incoming HTTP payloads
func ThemeModelFromMap(storeID primitive.ObjectID, payload map[string]interface{}) *models.Theme {
	t := &models.Theme{}
	t.StoreID = storeID
	if name, ok := payload["name"].(string); ok {
		t.Name = name
	}
	if desc, ok := payload["description"].(string); ok {
		t.Description = desc
	}
	if tokens, ok := payload["tokens"].(map[string]interface{}); ok {
		t.Tokens = tokens
	}
	if defaults, ok := payload["componentDefaults"].(map[string]interface{}); ok {
		t.ComponentDefaults = defaults
	}
	t.IsActive = false
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	return t
}
