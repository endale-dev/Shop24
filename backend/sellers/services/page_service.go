package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Endale2/DRPS/sellers/models"
	"github.com/Endale2/DRPS/sellers/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PageService struct {
	repo *repositories.PageRepository
}

func NewPageService() *PageService {
	return &PageService{repo: repositories.NewPageRepository()}
}

// Basic validation placeholder: ensure blocks have component keys and ids
func (s *PageService) validate(page *models.Page) error {
	if page.Path == "" {
		return errors.New("path required")
	}
	for _, b := range page.Blocks {
		if b.ComponentKey == "" || b.ID == "" {
			return errors.New("block id and componentKey required")
		}
	}
	return nil
}

func (s *PageService) Create(ctx context.Context, p *models.Page) error {
	if err := s.validate(p); err != nil {
		return err
	}
	return s.repo.Create(ctx, p)
}

func (s *PageService) Update(ctx context.Context, id primitive.ObjectID, p *models.Page) error {
	if err := s.validate(p); err != nil {
		return err
	}
	return s.repo.Update(ctx, id, p)
}

// Publish creates a snapshot (serialized blocks + tokens) and stores it on the page document
func (s *PageService) Publish(ctx context.Context, id primitive.ObjectID, p *models.Page) error {
	// validate
	if err := s.validate(p); err != nil {
		return err
	}

	// create a snapshot object for fast rendering; we include blocks and metadata
	snap := map[string]interface{}{
		"path":        p.Path,
		"title":       p.Title,
		"blocks":      p.Blocks,
		"publishedAt": time.Now(),
	}

	// ensure snapshot serializes (safety)
	if _, err := json.Marshal(snap); err != nil {
		return err
	}

	return s.repo.PublishSnapshot(ctx, id, snap)
}

func (s *PageService) GetPublishedByStoreAndPath(ctx context.Context, storeID primitive.ObjectID, path string) (*models.Page, error) {
	p, err := s.repo.FindByStoreAndPath(ctx, storeID, path)
	if err != nil {
		return nil, err
	}
	if p.Status != "published" {
		return nil, errors.New("page not published")
	}
	return p, nil
}

// ListByStore returns all pages for a store (draft and published)
func (s *PageService) ListByStore(ctx context.Context, storeID primitive.ObjectID) ([]*models.Page, error) {
	return s.repo.ListByStore(ctx, storeID)
}

// GetByID returns a single page by its ID
func (s *PageService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Page, error) {
	return s.repo.FindByID(ctx, id)
}

// Delete removes a page by its ID
func (s *PageService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteByID(ctx, id)
}

// CreateDefaultHomePage creates and publishes a default home page if one doesn't exist
func (s *PageService) CreateDefaultHomePage(ctx context.Context, storeID primitive.ObjectID) (*models.Page, error) {
	// Check if already exists
	existing, err := s.repo.FindByStoreAndPath(ctx, storeID, "/")
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new default page
	now := time.Now()
	page := &models.Page{
		ID:         primitive.NewObjectID(),
		StoreID:    storeID,
		Path:       "/",
		Title:      "Home",
		Visibility: "public",
		Status:     "published",
		Blocks: []models.PageBlock{
			{
				ID:           "default-hero",
				ComponentKey: "hero",
				Props: map[string]interface{}{
					"title":      "Welcome to our Store",
					"subtitle":   "We offer the best products for you.",
					"align":      "center",
					"background": "primary",
					"showButton": true,
					"buttonText": "Shop Now",
					"buttonLink": "/products",
				},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create Snapshot immediately
	page.Snapshot = map[string]interface{}{
		"path":        page.Path,
		"title":       page.Title,
		"blocks":      page.Blocks,
		"publishedAt": now,
	}

	if err := s.repo.Create(ctx, page); err != nil {
		return nil, err
	}

	return page, nil
}

// PageModelFromJSON is a lightweight payload type for incoming JSON -> model conversion
type PageModelFromJSON struct {
	Path       string                 `json:"path"`
	Title      string                 `json:"title"`
	SEO        map[string]interface{} `json:"seo"`
	Visibility string                 `json:"visibility"`
	Status     string                 `json:"status"`
	Blocks     []models.PageBlock     `json:"blocks"`
}

func (p *PageModelFromJSON) ToModel(storeID primitive.ObjectID) *models.Page {
	m := &models.Page{
		StoreID:    storeID,
		Path:       p.Path,
		Title:      p.Title,
		SEO:        p.SEO,
		Visibility: p.Visibility,
		Status:     p.Status,
		Blocks:     p.Blocks,
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return m
}
