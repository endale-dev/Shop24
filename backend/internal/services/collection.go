package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/domain/input"
	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

type CollectionService struct {
	stores      interfaces.StoreRepository
	products    interfaces.ProductRepository
	collections interfaces.CollectionRepository
}

func NewCollectionService(stores interfaces.StoreRepository, products interfaces.ProductRepository, collections interfaces.CollectionRepository) *CollectionService {
	return &CollectionService{stores: stores, products: products, collections: collections}
}

type CollectionInput struct {
	Name        string
	Slug        string
	Description string
	ImageURL    string
	ProductIDs  []string
	Published   bool
}

func (s *CollectionService) Create(ctx context.Context, ownerID string, in CollectionInput) (models.Collection, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Collection{}, err
	}
	collection, err := s.prepare(ctx, store.ID, models.Collection{StoreID: store.ID}, in)
	if err != nil {
		return models.Collection{}, err
	}
	now := time.Now().UTC()
	collection.CreatedAt = now
	collection.UpdatedAt = now
	created, err := s.collections.Create(ctx, collection)
	if err != nil {
		if appErr, ok := conflict(err); ok {
			return models.Collection{}, appErr
		}
		return models.Collection{}, err
	}
	return created, nil
}

func (s *CollectionService) Update(ctx context.Context, ownerID, id string, in CollectionInput) (models.Collection, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Collection{}, err
	}
	current, err := s.collections.FindByID(ctx, store.ID, id)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Collection{}, apperr.New(http.StatusNotFound, "not_found", "Collection not found")
	}
	if err != nil {
		return models.Collection{}, err
	}
	collection, err := s.prepare(ctx, store.ID, current, in)
	if err != nil {
		return models.Collection{}, err
	}
	collection.UpdatedAt = time.Now().UTC()
	updated, err := s.collections.Update(ctx, collection)
	if err != nil {
		if appErr, ok := conflict(err); ok {
			return models.Collection{}, appErr
		}
		if errors.Is(err, interfaces.ErrNotFound) {
			return models.Collection{}, apperr.New(http.StatusNotFound, "not_found", "Collection not found")
		}
		return models.Collection{}, err
	}
	return updated, nil
}

func (s *CollectionService) Delete(ctx context.Context, ownerID, id string) error {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return err
	}
	err = s.collections.Delete(ctx, store.ID, id)
	if errors.Is(err, interfaces.ErrNotFound) {
		return apperr.New(http.StatusNotFound, "not_found", "Collection not found")
	}
	return err
}

func (s *CollectionService) Get(ctx context.Context, ownerID, id string) (models.Collection, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Collection{}, err
	}
	collection, err := s.collections.FindByID(ctx, store.ID, id)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Collection{}, apperr.New(http.StatusNotFound, "not_found", "Collection not found")
	}
	return collection, err
}

func (s *CollectionService) List(ctx context.Context, ownerID string) ([]models.Collection, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return nil, err
	}
	items, err := s.collections.List(ctx, store.ID, interfaces.CollectionFilter{Limit: 100})
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.Collection{}
	}
	return items, nil
}

func (s *CollectionService) PublicList(ctx context.Context, slug string) ([]models.Collection, error) {
	store, err := publishedStore(ctx, s.stores, slug)
	if err != nil {
		return nil, err
	}
	published := true
	items, err := s.collections.List(ctx, store.ID, interfaces.CollectionFilter{Published: &published, Limit: 100})
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.Collection{}
	}
	return items, nil
}

func (s *CollectionService) PublicGet(ctx context.Context, slug, collectionSlug string) (models.Collection, []models.Product, models.PublicStore, error) {
	store, err := publishedStore(ctx, s.stores, slug)
	if err != nil {
		return models.Collection{}, nil, models.PublicStore{}, err
	}
	collection, err := s.collections.FindBySlug(ctx, store.ID, collectionSlug)
	if errors.Is(err, interfaces.ErrNotFound) || (err == nil && !collection.Published) {
		return models.Collection{}, nil, models.PublicStore{}, apperr.New(http.StatusNotFound, "not_found", "Collection not found")
	}
	if err != nil {
		return models.Collection{}, nil, models.PublicStore{}, err
	}
	products := make([]models.Product, 0, len(collection.ProductIDs))
	for _, id := range collection.ProductIDs {
		product, err := s.products.FindByID(ctx, store.ID, id)
		if errors.Is(err, interfaces.ErrNotFound) || (err == nil && !product.Published) {
			continue
		}
		if err != nil {
			return models.Collection{}, nil, models.PublicStore{}, err
		}
		products = append(products, presentProduct(product))
	}
	return collection, products, store.Public(), nil
}

func (s *CollectionService) prepare(ctx context.Context, storeID string, current models.Collection, in CollectionInput) (models.Collection, error) {
	name, err := input.Text(in.Name, "Collection name", 2, 80)
	if err != nil {
		return models.Collection{}, apperr.New(http.StatusBadRequest, "invalid_name", err.Error())
	}
	rawSlug := in.Slug
	if rawSlug == "" {
		rawSlug = name
	}
	slug, err := input.Slug(rawSlug)
	if err != nil {
		return models.Collection{}, apperr.New(http.StatusBadRequest, "invalid_slug", "Use 3 to 40 letters, numbers, and hyphens")
	}
	description, err := input.Optional(in.Description, 500)
	if err != nil {
		return models.Collection{}, apperr.New(http.StatusBadRequest, "invalid_description", "Description must be 500 characters or less")
	}
	image, err := input.URL(in.ImageURL)
	if err != nil {
		return models.Collection{}, apperr.New(http.StatusBadRequest, "invalid_images", "The collection image must start with http:// or https://")
	}
	ids, err := s.ownedProducts(ctx, storeID, in.ProductIDs)
	if err != nil {
		return models.Collection{}, err
	}
	current.Name = name
	current.Slug = slug
	current.Description = description
	current.ImageURL = image
	current.ProductIDs = ids
	current.Published = in.Published
	return current, nil
}

func (s *CollectionService) ownedProducts(ctx context.Context, storeID string, ids []string) ([]string, error) {
	if len(ids) > 100 {
		return nil, apperr.New(http.StatusBadRequest, "invalid_products", "A collection can hold at most 100 products")
	}
	out := make([]string, 0, len(ids))
	seen := map[string]struct{}{}
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		if _, err := s.products.FindByID(ctx, storeID, id); errors.Is(err, interfaces.ErrNotFound) {
			return nil, apperr.New(http.StatusBadRequest, "invalid_products", "Every product in a collection must belong to this store")
		} else if err != nil {
			return nil, err
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}
