package mongorepo

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	gomongo "go.mongodb.org/mongo-driver/mongo"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

func objectID(hex string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return primitive.NilObjectID, interfaces.ErrNotFound
	}
	return id, nil
}

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gomongo.ErrNoDocuments) {
		return interfaces.ErrNotFound
	}
	if gomongo.IsDuplicateKeyError(err) {
		return &interfaces.DuplicateKey{Field: duplicateField(err)}
	}
	return err
}

func duplicateField(err error) string {
	message := err.Error()
	switch {
	case strings.Contains(message, "store_sku"):
		return "sku"
	case strings.Contains(message, "ownerId"):
		return "ownerId"
	case strings.Contains(message, "collection_slug"):
		return "collection_slug"
	case strings.Contains(message, "slug"):
		return "slug"
	case strings.Contains(message, "googleSub"):
		return "googleSub"
	case strings.Contains(message, "email"):
		return "email"
	default:
		return ""
	}
}

func limitOf(value, fallback, max int) int64 {
	if value <= 0 {
		value = fallback
	}
	if value > max {
		value = max
	}
	return int64(value)
}

func utc(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}

type userDoc struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"passwordHash"`
	GoogleSub    string             `bson:"googleSub,omitempty"`
	ResetHash    string             `bson:"resetHash,omitempty"`
	ResetExpires time.Time          `bson:"resetExpires,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

func (d userDoc) model() models.User {
	return models.User{
		ID: d.ID.Hex(), Name: d.Name, Email: d.Email, PasswordHash: d.PasswordHash,
		GoogleSub: d.GoogleSub, ResetHash: d.ResetHash, ResetExpires: d.ResetExpires,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type variantOptionDoc struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}

type variantDoc struct {
	ID         string             `bson:"id"`
	Options    []variantOptionDoc `bson:"options"`
	SKU        string             `bson:"sku"`
	PriceCents int64              `bson:"priceCents"`
	Stock      int                `bson:"stock"`
	ImageURL   string             `bson:"imageUrl"`
}

type productDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	StoreID     string             `bson:"storeId"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	PriceCents  int64              `bson:"priceCents"`
	Images      []string           `bson:"images"`
	Category    string             `bson:"category"`
	OptionNames []string           `bson:"optionNames"`
	Variants    []variantDoc       `bson:"variants"`
	SKU         string             `bson:"sku"`
	Stock       int                `bson:"stock"`
	Published   bool               `bson:"published"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

func variantDocs(variants []models.Variant) []variantDoc {
	if variants == nil {
		return []variantDoc{}
	}
	out := make([]variantDoc, 0, len(variants))
	for _, variant := range variants {
		options := make([]variantOptionDoc, 0, len(variant.Options))
		for _, option := range variant.Options {
			options = append(options, variantOptionDoc{Name: option.Name, Value: option.Value})
		}
		out = append(out, variantDoc{
			ID: variant.ID, Options: options, SKU: variant.SKU,
			PriceCents: variant.PriceCents, Stock: variant.Stock, ImageURL: variant.ImageURL,
		})
	}
	return out
}

func variantModels(docs []variantDoc) []models.Variant {
	if docs == nil {
		return nil
	}
	out := make([]models.Variant, 0, len(docs))
	for _, doc := range docs {
		options := make([]models.VariantOption, 0, len(doc.Options))
		for _, option := range doc.Options {
			options = append(options, models.VariantOption{Name: option.Name, Value: option.Value})
		}
		out = append(out, models.Variant{
			ID: doc.ID, Options: options, SKU: doc.SKU,
			PriceCents: doc.PriceCents, Stock: doc.Stock, ImageURL: doc.ImageURL,
		})
	}
	return out
}

func newProductDoc(product models.Product) productDoc {
	images := product.Images
	if images == nil {
		images = []string{}
	}
	names := product.OptionNames
	if names == nil {
		names = []string{}
	}
	return productDoc{
		StoreID: product.StoreID, Name: product.Name, Description: product.Description,
		PriceCents: product.PriceCents, Images: images, Category: product.Category,
		OptionNames: names, Variants: variantDocs(product.Variants),
		SKU: product.SKU, Stock: product.Stock, Published: product.Published,
		CreatedAt: utc(product.CreatedAt), UpdatedAt: utc(product.UpdatedAt),
	}
}

func (d productDoc) model() models.Product {
	images := d.Images
	if images == nil {
		images = []string{}
	}
	names := d.OptionNames
	if names == nil {
		names = []string{}
	}
	return models.Product{
		ID: d.ID.Hex(), StoreID: d.StoreID, Name: d.Name, Description: d.Description,
		PriceCents: d.PriceCents, Images: images, Category: d.Category,
		OptionNames: names, Variants: variantModels(d.Variants),
		SKU: d.SKU, Stock: d.Stock, Published: d.Published, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type orderItemDoc struct {
	ProductID      string `bson:"productId"`
	VariantID      string `bson:"variantId"`
	Name           string `bson:"name"`
	OptionLabel    string `bson:"optionLabel"`
	SKU            string `bson:"sku"`
	UnitPriceCents int64  `bson:"unitPriceCents"`
	Quantity       int    `bson:"quantity"`
	LineTotalCents int64  `bson:"lineTotalCents"`
}

type customerDoc struct {
	Name    string `bson:"name"`
	Email   string `bson:"email"`
	Phone   string `bson:"phone"`
	Address string `bson:"address"`
	City    string `bson:"city"`
	Notes   string `bson:"notes"`
}

type orderDoc struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	StoreID       string             `bson:"storeId"`
	Status        string             `bson:"status"`
	Customer      customerDoc        `bson:"customer"`
	Items         []orderItemDoc     `bson:"items"`
	SubtotalCents int64              `bson:"subtotalCents"`
	TotalCents    int64              `bson:"totalCents"`
	PaymentMethod string             `bson:"paymentMethod"`
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"`
}

func newOrderDoc(order models.Order) orderDoc {
	items := make([]orderItemDoc, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, orderItemDoc{
			ProductID: item.ProductID, VariantID: item.VariantID, Name: item.Name, OptionLabel: item.OptionLabel, SKU: item.SKU,
			UnitPriceCents: item.UnitPriceCents, Quantity: item.Quantity, LineTotalCents: item.LineTotalCents,
		})
	}
	return orderDoc{
		StoreID: order.StoreID, Status: string(order.Status),
		Customer: customerDoc{
			Name: order.Customer.Name, Email: order.Customer.Email, Phone: order.Customer.Phone,
			Address: order.Customer.Address, City: order.Customer.City, Notes: order.Customer.Notes,
		},
		Items: items, SubtotalCents: order.SubtotalCents, TotalCents: order.TotalCents,
		PaymentMethod: order.PaymentMethod, CreatedAt: utc(order.CreatedAt), UpdatedAt: utc(order.UpdatedAt),
	}
}

func (d orderDoc) model() models.Order {
	items := make([]models.OrderItem, 0, len(d.Items))
	for _, item := range d.Items {
		items = append(items, models.OrderItem{
			ProductID: item.ProductID, VariantID: item.VariantID, Name: item.Name, OptionLabel: item.OptionLabel, SKU: item.SKU,
			UnitPriceCents: item.UnitPriceCents, Quantity: item.Quantity, LineTotalCents: item.LineTotalCents,
		})
	}
	return models.Order{
		ID: d.ID.Hex(), StoreID: d.StoreID, Status: models.OrderStatus(d.Status),
		Customer: models.Customer{
			Name: d.Customer.Name, Email: d.Customer.Email, Phone: d.Customer.Phone,
			Address: d.Customer.Address, City: d.Customer.City, Notes: d.Customer.Notes,
		},
		Items: items, SubtotalCents: d.SubtotalCents, TotalCents: d.TotalCents,
		PaymentMethod: d.PaymentMethod, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}
