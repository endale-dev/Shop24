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
	case strings.Contains(message, "slug"):
		return "slug"
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
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

func (d userDoc) model() models.User {
	return models.User{
		ID: d.ID.Hex(), Name: d.Name, Email: d.Email, PasswordHash: d.PasswordHash,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type productDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	StoreID     string             `bson:"storeId"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	PriceCents  int64              `bson:"priceCents"`
	Images      []string           `bson:"images"`
	Category    string             `bson:"category"`
	SKU         string             `bson:"sku"`
	Stock       int                `bson:"stock"`
	Published   bool               `bson:"published"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

func newProductDoc(product models.Product) productDoc {
	images := product.Images
	if images == nil {
		images = []string{}
	}
	return productDoc{
		StoreID: product.StoreID, Name: product.Name, Description: product.Description,
		PriceCents: product.PriceCents, Images: images, Category: product.Category,
		SKU: product.SKU, Stock: product.Stock, Published: product.Published,
		CreatedAt: utc(product.CreatedAt), UpdatedAt: utc(product.UpdatedAt),
	}
}

func (d productDoc) model() models.Product {
	images := d.Images
	if images == nil {
		images = []string{}
	}
	return models.Product{
		ID: d.ID.Hex(), StoreID: d.StoreID, Name: d.Name, Description: d.Description,
		PriceCents: d.PriceCents, Images: images, Category: d.Category, SKU: d.SKU,
		Stock: d.Stock, Published: d.Published, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

type orderItemDoc struct {
	ProductID      string `bson:"productId"`
	Name           string `bson:"name"`
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
			ProductID: item.ProductID, Name: item.Name, SKU: item.SKU,
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
			ProductID: item.ProductID, Name: item.Name, SKU: item.SKU,
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
