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
