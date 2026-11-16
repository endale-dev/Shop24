package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	stores      interfaces.StoreRepository
	products    interfaces.ProductRepository
	collections interfaces.CollectionRepository
}

func NewProductService(stores interfaces.StoreRepository, products interfaces.ProductRepository, collections interfaces.CollectionRepository) *ProductService {
	return &ProductService{stores: stores, products: products, collections: collections}
}

type VariantInput struct {
	ID         string
	Options    []models.VariantOption
	SKU        string
	PriceCents int64
	Stock      int
	ImageURL   string
}

type ProductInput struct {
	Name        string
	Description string
	Images      []string
	Category    string
	OptionNames []string
	Variants    []VariantInput
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
	if err := s.ensureSKUs(ctx, store.ID, "", product.Variants); err != nil {
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
	kept := map[string]int{}
	for _, variant := range current.WithVariants().Variants {
		kept[variant.ID] = variant.Stock
	}
	product, err := validateProduct(current, in)
	if err != nil {
		return models.Product{}, err
	}
	for i, variant := range product.Variants {
		if stock, ok := kept[variant.ID]; ok {
			product.Variants[i].Stock = stock
		}
	}
	product = product.WithVariants()
	if err := s.ensureSKUs(ctx, store.ID, product.ID, product.Variants); err != nil {
		return models.Product{}, err
	}
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

func (s *ProductService) SetStock(ctx context.Context, ownerID, productID, variantID string, stock int) (models.Product, error) {
	if stock < 0 || stock > 1_000_000 {
		return models.Product{}, apperr.New(http.StatusBadRequest, "invalid_stock", "Stock must be between 0 and 1000000")
	}
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
	current = current.WithVariants()
	if variantID == "" {
		if len(current.Variants) != 1 {
			return models.Product{}, apperr.New(http.StatusBadRequest, "choose_variant", "Choose which version to restock")
		}
		variantID = current.Variants[0].ID
	} else if !hasVariant(current, variantID) {
		return models.Product{}, apperr.New(http.StatusNotFound, "not_found", "That version no longer exists")
	}
	if err := s.products.SetStock(ctx, store.ID, productID, variantID, stock, time.Now().UTC()); err != nil {
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
	if err != nil {
		return err
	}
	return s.collections.DetachProduct(ctx, store.ID, productID)
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
	if len(filter.Query) > 80 {
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

func (s *ProductService) ensureSKUs(ctx context.Context, storeID, exceptID string, variants []models.Variant) error {
	seen := map[string]struct{}{}
	for _, variant := range variants {
		if _, ok := seen[variant.SKU]; ok {
			return apperr.New(http.StatusConflict, "sku_taken", "Each version needs its own SKU")
		}
		seen[variant.SKU] = struct{}{}
		used, err := s.products.SKUInUse(ctx, storeID, variant.SKU, exceptID)
		if err != nil {
			return err
		}
		if used {
			return apperr.New(http.StatusConflict, "sku_taken", "A product with this SKU already exists")
		}
	}
	return nil
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
	names, err := cleanOptionNames(in.OptionNames)
	if err != nil {
		return models.Product{}, err
	}
	variants, err := cleanVariants(current, names, in.Variants)
	if err != nil {
		return models.Product{}, err
	}
	images, err := cleanImages(in.Images)
	if err != nil {
		return models.Product{}, err
	}
	current.Name = name
	current.Description = description
	current.Category = category
	current.OptionNames = names
	current.Variants = variants
	current.Images = images
	current.Published = in.Published
	return current.WithVariants(), nil
}

func cleanOptionNames(names []string) ([]string, error) {
	if len(names) > 3 {
		return nil, apperr.New(http.StatusBadRequest, "invalid_options", "A product can have at most 3 options")
	}
	out := make([]string, 0, len(names))
	seen := map[string]struct{}{}
	for _, name := range names {
		value, err := input.Text(name, "Option name", 1, 24)
		if err != nil {
			return nil, apperr.New(http.StatusBadRequest, "invalid_options", "Option names must be 1 to 24 characters")
		}
		key := strings.ToLower(value)
		if _, ok := seen[key]; ok {
			return nil, apperr.New(http.StatusBadRequest, "invalid_options", "Option names must be different")
		}
		seen[key] = struct{}{}
		out = append(out, value)
	}
	return out, nil
}

func cleanVariants(current models.Product, names []string, in []VariantInput) ([]models.Variant, error) {
	if len(in) == 0 || len(in) > 30 {
		return nil, apperr.New(http.StatusBadRequest, "invalid_variants", "Add between 1 and 30 versions")
	}
	if len(names) == 0 && len(in) != 1 {
		return nil, apperr.New(http.StatusBadRequest, "invalid_variants", "Name an option before adding another version")
	}
	known := map[string]struct{}{}
	for _, variant := range current.WithVariants().Variants {
		if variant.ID != "" {
			known[variant.ID] = struct{}{}
		}
	}
	out := make([]models.Variant, 0, len(in))
	combos := map[string]struct{}{}
	skus := map[string]struct{}{}
	for _, item := range in {
		options := make([]models.VariantOption, 0, len(names))
		combo := make([]string, 0, len(names))
		for _, name := range names {
			raw := ""
			for _, option := range item.Options {
				if strings.EqualFold(strings.TrimSpace(option.Name), name) {
					raw = option.Value
					break
				}
			}
			value, err := input.Text(raw, name, 1, 40)
			if err != nil {
				return nil, apperr.New(http.StatusBadRequest, "invalid_variants", name+" needs a value on every version")
			}
			options = append(options, models.VariantOption{Name: name, Value: value})
			combo = append(combo, strings.ToLower(value))
		}
		key := strings.Join(combo, "\x00")
		if _, ok := combos[key]; ok {
			return nil, apperr.New(http.StatusBadRequest, "invalid_variants", "Two versions cannot have the same options")
		}
		combos[key] = struct{}{}
		sku := strings.TrimSpace(item.SKU)
		if !skuPattern.MatchString(sku) {
			return nil, apperr.New(http.StatusBadRequest, "invalid_sku", "SKU must be 1 to 40 letters, numbers, hyphens, or underscores")
		}
		if _, ok := skus[sku]; ok {
			return nil, apperr.New(http.StatusConflict, "sku_taken", "Each version needs its own SKU")
		}
		skus[sku] = struct{}{}
		if item.PriceCents < 1 || item.PriceCents > 100_000_000 {
			return nil, apperr.New(http.StatusBadRequest, "invalid_price", "Enter a price between 0.01 and 1000000")
		}
		if item.Stock < 0 || item.Stock > 1_000_000 {
			return nil, apperr.New(http.StatusBadRequest, "invalid_stock", "Stock must be between 0 and 1000000")
		}
		image, err := input.URL(item.ImageURL)
		if err != nil {
			return nil, apperr.New(http.StatusBadRequest, "invalid_images", "Version image URLs must start with http:// or https://")
		}
		id := strings.TrimSpace(item.ID)
		if _, ok := known[id]; !ok {
			id = newVariantID()
		}
		out = append(out, models.Variant{
			ID: id, Options: options, SKU: sku, PriceCents: item.PriceCents, Stock: item.Stock, ImageURL: image,
		})
	}
	return out, nil
}

func hasVariant(product models.Product, id string) bool {
	for _, variant := range product.Variants {
		if variant.ID == id {
			return true
		}
	}
	return false
}

func newVariantID() string {
	var buf [12]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))[:24]
	}
	return hex.EncodeToString(buf[:])
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
	product = product.WithVariants()
	if product.Images == nil {
		product.Images = []string{}
	}
	product.LowStock = false
	for _, variant := range product.Variants {
		if variant.Stock <= LowStockThreshold {
			product.LowStock = true
			break
		}
	}
	return product
}

func presentProducts(products []models.Product) []models.Product {
	out := make([]models.Product, 0, len(products))
	for _, product := range products {
		out = append(out, presentProduct(product))
	}
	return out
}
