package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/domain/models"
	"shop24/internal/domain/storefront"
	"shop24/internal/repositories/interfaces"
)

type StorefrontService struct {
	stores interfaces.StoreRepository
}

func NewStorefrontService(stores interfaces.StoreRepository) *StorefrontService {
	return &StorefrontService{stores: stores}
}

func (s *StorefrontService) Save(ctx context.Context, ownerID string, cfg models.StorefrontConfig) (models.Store, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Store{}, err
	}
	cfg = storefront.Normalize(cfg)
	now := time.Now().UTC()
	if err := s.stores.UpdateDraft(ctx, ownerID, store.ID, cfg, nil, now); err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return models.Store{}, apperr.New(http.StatusNotFound, "not_found", "Store not found")
		}
		return models.Store{}, err
	}
	store.Draft = cfg
	store.UpdatedAt = now
	return store, nil
}

func (s *StorefrontService) Publish(ctx context.Context, ownerID string) (models.Store, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Store{}, err
	}
	draft := storefront.Normalize(store.Draft)
	now := time.Now().UTC()
	if err := s.stores.UpdateDraft(ctx, ownerID, store.ID, draft, nil, now); err != nil {
		return models.Store{}, err
	}
	if err := s.stores.PublishDraft(ctx, ownerID, store.ID, now); err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return models.Store{}, apperr.New(http.StatusNotFound, "not_found", "Store not found")
		}
		return models.Store{}, err
	}
	store.Draft = draft
	store.Storefront = draft
	store.UpdatedAt = now
	return store, nil
}
