package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/domain/input"
	"shop24/internal/domain/models"
	"shop24/internal/domain/storefront"
	"shop24/internal/repositories/interfaces"
)

var currencies = map[string]struct{}{"USD": {}, "EUR": {}, "GBP": {}}

type StoreService struct {
	stores interfaces.StoreRepository
}

func NewStoreService(stores interfaces.StoreRepository) *StoreService {
	return &StoreService{stores: stores}
}

type StoreInput struct {
	Name        string
	Slug        string
	Description string
	LogoURL     string
	Currency    string
}

func (s *StoreService) Create(ctx context.Context, ownerID string, in StoreInput) (models.Store, error) {
	if _, err := s.stores.FindByOwner(ctx, ownerID); err == nil {
		return models.Store{}, apperr.New(http.StatusConflict, "store_exists", "You already have a store")
	} else if !errors.Is(err, interfaces.ErrNotFound) {
		return models.Store{}, err
	}
	store, err := s.prepare(models.Store{OwnerID: ownerID, Currency: "USD"}, in, false)
	if err != nil {
		return models.Store{}, err
	}
	design := storefront.Default(store.Name, store.Description)
	now := time.Now().UTC()
	store.Storefront = design
	store.Draft = design
	store.Published = false
	store.CreatedAt = now
	store.UpdatedAt = now
	created, err := s.stores.Create(ctx, store)
	if err != nil {
		if appErr, ok := conflict(err); ok {
			return models.Store{}, appErr
		}
		return models.Store{}, err
	}
	return created, nil
}

func (s *StoreService) Update(ctx context.Context, ownerID string, in StoreInput) (models.Store, error) {
	current, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Store{}, err
	}
	updated, err := s.prepare(current, in, true)
	if err != nil {
		return models.Store{}, err
	}
	updated.UpdatedAt = time.Now().UTC()
	if err := s.stores.UpdateProfile(ctx, updated); err != nil {
		if appErr, ok := conflict(err); ok {
			return models.Store{}, appErr
		}
		if errors.Is(err, interfaces.ErrNotFound) {
			return models.Store{}, apperr.ErrNotFound
		}
		return models.Store{}, err
	}
	return updated, nil
}

func (s *StoreService) SetPublished(ctx context.Context, ownerID string, published bool) (models.Store, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Store{}, err
	}
	now := time.Now().UTC()
	if err := s.stores.SetPublished(ctx, ownerID, store.ID, published, now); err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return models.Store{}, apperr.ErrNotFound
		}
		return models.Store{}, err
	}
	store.Published = published
	store.UpdatedAt = now
	return store, nil
}

func (s *StoreService) Mine(ctx context.Context, ownerID string) (models.Store, error) {
	store, err := s.stores.FindByOwner(ctx, ownerID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Store{}, apperr.New(http.StatusNotFound, "not_found", "You have not created a store yet")
	}
	return store, err
}

func (s *StoreService) Public(ctx context.Context, slug string) (models.PublicStore, error) {
	store, err := publishedStore(ctx, s.stores, slug)
	if err != nil {
		return models.PublicStore{}, err
	}
	return store.Public(), nil
}

func (s *StoreService) prepare(current models.Store, in StoreInput, keepCurrency bool) (models.Store, error) {
	name, err := input.Text(in.Name, "Store name", 2, 60)
	if err != nil {
		return models.Store{}, apperr.New(http.StatusBadRequest, "invalid_name", err.Error())
	}
	slug, err := input.Slug(in.Slug)
	if err != nil {
		return models.Store{}, apperr.New(http.StatusBadRequest, "invalid_slug", err.Error())
	}
	description, err := input.Optional(in.Description, 500)
	if err != nil {
		return models.Store{}, apperr.New(http.StatusBadRequest, "invalid_description", "Description must be 500 characters or less")
	}
	logo, err := input.URL(in.LogoURL)
	if err != nil {
		return models.Store{}, apperr.New(http.StatusBadRequest, "invalid_logo", err.Error())
	}
	currency := current.Currency
	if !keepCurrency || in.Currency != "" {
		currency = in.Currency
		if currency == "" {
			currency = "USD"
		}
		if _, ok := currencies[currency]; !ok {
			return models.Store{}, apperr.New(http.StatusBadRequest, "invalid_currency", "Currency must be USD, EUR, or GBP")
		}
	}
	current.Name = name
	current.Slug = slug
	current.Description = description
	current.LogoURL = logo
	current.Currency = currency
	return current, nil
}
