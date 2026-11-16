package services

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"shop24/internal/apperr"
	"shop24/internal/domain/input"
	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

const LowStockThreshold = 5

var skuPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,39}$`)

type ProductService struct {
	stores   interfaces.StoreRepository
	products interfaces.ProductRepository
}

func NewProductService(stores interfaces.StoreRepository, products interfaces.ProductRepository) *ProductService {
	return &ProductService{stores: stores, products: products}
}

type ProductInput struct {
	Name        string
	Description string
	PriceCents  int64
	Images      []string
	Category    string
	SKU         string
	Stock       int
	Published   bool
}

func (s *ProductService) Create(ctx context.Context, ownerID string, in ProductInput) (models.Product, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Product{}, err
	}
	product, err := validateProduct(models.Product{StoreID: store.ID}, in)
	if err != nil {
		return models.Product{}, err
	}
	now := time.Now().UTC()
	product.CreatedAt = now
	product.UpdatedAt = now
	created, err := s.products.Create(ctx, product)
	if err != nil {
		if appErr, ok := conflict(err); ok {
			return models.Product{}, appErr
		}
		return models.Product{}, err
	}
	return presentProduct(created), nil
}

func (s *ProductService) Update(ctx context.Context, ownerID, productID string, in ProductInput) (models.Product, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Product{}, err
	}
	current, err := s.products.FindByID(ctx, store.ID, productID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Product{}, apperr.New(http.StatusNotFound, "not_found", "Product not found")
	}
	if err != nil {
		return models.Product{}, err
	}
	keptStock := current.Stock
	product, err := validateProduct(current, in)
	if err != nil {
		return models.Product{}, err
	}
	product.Stock = keptStock
	product.UpdatedAt = time.Now().UTC()
	updated, err := s.products.Update(ctx, product)
	if err != nil {
		if appErr, ok := conflict(err); ok {
			return models.Product{}, appErr
		}
		if errors.Is(err, interfaces.ErrNotFound) {
			return models.Product{}, apperr.New(http.StatusNotFound, "not_found", "Product not found")
		}
		return models.Product{}, err
	}
	return presentProduct(updated), nil
}

func (s *ProductService) SetStock(ctx context.Context, ownerID, productID string, stock int) (models.Product, error) {
	if stock < 0 || stock > 1_000_000 {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_stock", "Stock must be between 0 and 1000000")
	}
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Product{}, err
	}
	if err := s.products.SetStock(ctx, store.ID, productID, stock, time.Now().UTC()); err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return models.Product{}, apperr.New(http.StatusNotFound, "not_found", "Product not found")
		}
		return models.Product{}, err
	}
	product, err := s.products.FindByID(ctx, store.ID, productID)
	if err != nil {
		return models.Product{}, err
	}
	return presentProduct(product), nil
}

func (s *ProductService) Delete(ctx context.Context, ownerID, productID string) error {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return err
	}
	err = s.products.Delete(ctx, store.ID, productID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return apperr.New(http.StatusNotFound, "not_found", "Product not found")
	}
	return err
}

func (s *ProductService) Get(ctx context.Context, ownerID, productID string) (models.Product, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Product{}, err
	}
	product, err := s.products.FindByID(ctx, store.ID, productID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Product{}, apperr.New(http.StatusNotFound, "not_found", "Product not found")
	}
	if err != nil {
		return models.Product{}, err
	}
	return presentProduct(product), nil
}

func (s *ProductService) List(ctx context.Context, ownerID string, filter interfaces.ProductFilter) ([]models.Product, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return nil, err
	}
	if len(filter.Query) > 80 || len(filter.Category) > 40 {
		return nil, apperr.New(http.StatusBadRequest, "invalid_filter", "Search filters are too long")
	}
	products, err := s.products.List(ctx, store.ID, filter)
	if err != nil {
		return nil, err
	}
	return presentProducts(products), nil
}

func (s *ProductService) PublicList(ctx context.Context, slug string) ([]models.Product, error) {
	store, err := publishedStore(ctx, s.stores, slug)
	if err != nil {
		return nil, err
	}
	published := true
	products, err := s.products.List(ctx, store.ID, interfaces.ProductFilter{Published: &published, Limit: 100})
	if err != nil {
		return nil, err
	}
	return presentProducts(products), nil
}

func (s *ProductService) PublicGet(ctx context.Context, slug, productID string) (models.Product, models.PublicStore, error) {
	store, err := publishedStore(ctx, s.stores, slug)
	if err != nil {
		return models.Product{}, models.PublicStore{}, err
	}
	product, err := s.products.FindByID(ctx, store.ID, productID)
	if errors.Is(err, interfaces.ErrNotFound) || (err == nil && !product.Published) {
		return models.Product{}, models.PublicStore{}, apperr.New(http.StatusNotFound, "not_found", "Product not found")
	}
	if err != nil {
		return models.Product{}, models.PublicStore{}, err
	}
	return presentProduct(product), store.Public(), nil
}

func validateProduct(current models.Product, in ProductInput) (models.Product, error) {
	name, err := input.Text(in.Name, "Product name", 2, 120)
	if err != nil {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_name", err.Error())
	}
	description, err := input.Optional(in.Description, 4000)
	if err != nil {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_description", "Description must be 4000 characters or less")
	}
	category, err := input.Optional(in.Category, 40)
	if err != nil {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_category", "Category must be 40 characters or less")
	}
	sku := strings.TrimSpace(in.SKU)
	if !skuPattern.MatchString(sku) {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_sku", "SKU must be 1 to 40 letters, numbers, hyphens, or underscores")
	}
	if in.PriceCents < 1 || in.PriceCents > 100_000_000 {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_price", "Enter a price between 0.01 and 1000000")
	}
	if in.Stock < 0 || in.Stock > 1_000_000 {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_stock", "Stock must be between 0 and 1000000")
	}
	images, err := cleanImages(in.Images)
	if err != nil {
		return models.Product{}, err
	}
	current.Name = name
	current.Description = description
	current.Category = category
	current.SKU = sku
	current.PriceCents = in.PriceCents
	current.Stock = in.Stock
	current.Images = images
	current.Published = in.Published
	return current, nil
}

func cleanImages(images []string) ([]string, error) {
	if len(images) > 8 {
		return nil, apperr.New(http.StatusBadRequest, "invalid_images", "A product can have at most 8 images")
	}
	out := make([]string, 0, len(images))
	for _, image := range images {
		url, err := input.URL(image)
		if err != nil || (url == "" && image != "") {
			return nil, apperr.New(http.StatusBadRequest, "invalid_images", "Image URLs must start with http:// or https://")
		}
		if url != "" {
			out = append(out, url)
		}
	}
	return out, nil
}

func presentProduct(product models.Product) models.Product {
	if product.Images == nil {
		product.Images = []string{}
	}
	product.LowStock = product.Stock <= LowStockThreshold
	return product
}

func presentProducts(products []models.Product) []models.Product {
	out := make([]models.Product, 0, len(products))
	for _, product := range products {
		out = append(out, presentProduct(product))
	}
	return out
}
