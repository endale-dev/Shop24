package services

import (
	"context"
	"errors"
	"net/http"

	"shop24/internal/apperr"
	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

func conflict(err error) (error, bool) {
	var duplicate *interfaces.DuplicateKey
	if !errors.As(err, &duplicate) {
		return nil, false
	}
	switch duplicate.Field {
	case "email":
		return apperr.New(http.StatusConflict, "email_taken", "An account with this email already exists"), true
	case "slug":
		return apperr.New(http.StatusConflict, "slug_taken", "That store address is already taken"), true
	case "ownerId":
		return apperr.New(http.StatusConflict, "store_exists", "You already have a store"), true
	case "sku":
		return apperr.New(http.StatusConflict, "sku_taken", "A product with this SKU already exists"), true
	case "collection_slug":
		return apperr.New(http.StatusConflict, "slug_taken", "A collection with this address already exists"), true
	default:
		return apperr.New(http.StatusConflict, "conflict", "That value is already in use"), true
	}
}

func requireStore(ctx context.Context, stores interfaces.StoreRepository, ownerID string) (models.Store, error) {
	store, err := stores.FindByOwner(ctx, ownerID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Store{}, apperr.New(http.StatusConflict, "store_required", "Create a store before doing this")
	}
	if err != nil {
		return models.Store{}, err
	}
	return store, nil
}

func publishedStore(ctx context.Context, stores interfaces.StoreRepository, slug string) (models.Store, error) {
	store, err := stores.FindBySlug(ctx, slug)
	if errors.Is(err, interfaces.ErrNotFound) || (err == nil && !store.Published) {
		return models.Store{}, apperr.New(http.StatusNotFound, "not_found", "Store not found")
	}
	if err != nil {
		return models.Store{}, err
	}
	return store, nil
}
