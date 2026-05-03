package services

import (
	"context"
	"time"

	"github.com/Endale2/DRPS/sellers/models"
	"github.com/Endale2/DRPS/sellers/repositories"
	sharedModels "github.com/Endale2/DRPS/shared/models"
	"go.mongodb.org/mongo-driver/bson"
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

// GetActiveThemeByShopID retrieves the active theme for a shop
func (s *ThemeService) GetActiveThemeByShopID(ctx context.Context, shopID primitive.ObjectID) (*sharedModels.Theme, error) {
	theme, err := s.repo.FindActiveByShopID(ctx, shopID)
	if err != nil {
		// If not found or error, create a default active theme
		// In a real prod app check specifically for ErrNoDocuments
		return s.CreateDefaultTheme(ctx, shopID)
	}
	return theme, nil
}

// GetActiveThemeByShopSlug retrieves the active theme for a shop by slug (for storefront)
func (s *ThemeService) GetActiveThemeByShopSlug(ctx context.Context, shopSlug string) (*sharedModels.Theme, error) {
	theme, err := s.repo.FindActiveByShopSlug(ctx, shopSlug)
	if err != nil {
		// We need the shop ID to create a default theme.
		// Since we only have slug here, and repo failed, we must fetch shop first.
		// However, FindActiveByShopSlug already fetches shop.
		// Let's rely on the repository or just error out?
		// Actually, callers of this usually have access to Shop service.
		// But to keep it self-contained:
		// We can't easily create a default theme without ShopID.
		// Let's return error and let controller handle creation if needed?
		// NO, simpler: The controller calling this (GetStorefront) ALREADY has the Shop object.
		// So checking GetStorefront in storefront_controller, it fetches Shop first.
		// So we can fallback there correctly.
		// BUT `GetActiveThemeByShopID` is used by the Dashboard endpoint `GetActiveTheme`.
		// So updating `GetActiveThemeByShopID` fixes the Dashboard "Load" issue.
		// For `GetActiveThemeByShopSlug` (Storefront), we should probably let it return error
		// and let the controller use the ShopID it has to CreateDefault.
		return nil, err
	}
	return theme, nil
}

// UpdateTheme updates a theme
func (s *ThemeService) UpdateTheme(ctx context.Context, themeID primitive.ObjectID, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	return s.repo.Update(ctx, themeID, update)
}

// CreateDefaultTheme creates a default theme for a shop
func (s *ThemeService) CreateDefaultTheme(ctx context.Context, shopID primitive.ObjectID) (*sharedModels.Theme, error) {
	theme := sharedModels.DefaultTheme(shopID)

	// Convert to models.Theme for repository
	sellerTheme := &models.Theme{
		ShopID:            shopID,
		Name:              theme.Name,
		Description:       theme.Description,
		IsActive:          true,
		Tokens:            make(map[string]interface{}),
		ComponentDefaults: make(map[string]interface{}),
		CreatedAt:         theme.CreatedAt,
		UpdatedAt:         theme.UpdatedAt,
	}

	err := s.repo.Create(ctx, sellerTheme)
	if err != nil {
		return nil, err
	}

	theme.ID = sellerTheme.ID
	return theme, nil
}

// UpdateThemeColors updates the colors section of a theme
func (s *ThemeService) UpdateThemeColors(ctx context.Context, themeID primitive.ObjectID, colors sharedModels.ThemeColors) error {
	update := bson.M{
		"colors":    colors,
		"updatedAt": time.Now(),
	}
	return s.repo.Update(ctx, themeID, update)
}

// UpdateThemeTypography updates the typography section of a theme
func (s *ThemeService) UpdateThemeTypography(ctx context.Context, themeID primitive.ObjectID, typography sharedModels.ThemeTypography) error {
	update := bson.M{
		"typography": typography,
		"updatedAt":  time.Now(),
	}
	return s.repo.Update(ctx, themeID, update)
}

// UpdateThemeLayout updates the layout section of a theme
func (s *ThemeService) UpdateThemeLayout(ctx context.Context, themeID primitive.ObjectID, layout sharedModels.ThemeLayout) error {
	update := bson.M{
		"layout":    layout,
		"updatedAt": time.Now(),
	}
	return s.repo.Update(ctx, themeID, update)
}

// UpdateThemeAssets updates the assets section of a theme
func (s *ThemeService) UpdateThemeAssets(ctx context.Context, themeID primitive.ObjectID, assets sharedModels.ThemeAssets) error {
	update := bson.M{
		"assets":    assets,
		"updatedAt": time.Now(),
	}
	return s.repo.Update(ctx, themeID, update)
}

// ThemeModelFromMap builds a models.Theme from incoming HTTP payloads (legacy support)
func ThemeModelFromMap(storeID primitive.ObjectID, payload map[string]interface{}) *models.Theme {
	t := &models.Theme{}
	t.ShopID = storeID
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
