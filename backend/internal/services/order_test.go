package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

func TestPlaceOrderSnapshotsPriceAndRollsBackStock(t *testing.T) {
	stores := &memStores{items: []models.Store{{ID: "store-1", OwnerID: "user-1", Slug: "north", Published: true}}}
	products := &memProducts{items: map[string]models.Product{
		"coat": {ID: "coat", StoreID: "store-1", Name: "Coat", SKU: "COAT", PriceCents: 2500, Stock: 2, Published: true},
	}}
	orders := &memOrders{fail: true}
	service := NewOrderService(stores, products, orders)
	_, err := service.Place(context.Background(), "north", PlaceOrderInput{
		Customer: models.Customer{Name: "Ada Lovelace", Email: "ada@example.com", Address: "12 Analytical", City: "London"},
		Items:    []ItemInput{{ProductID: "coat", Quantity: 2, UnitPriceCents: 2500}},
	})
	if err == nil {
		t.Fatal("expected create failure")
	}
	if products.items["coat"].Stock != 2 {
		t.Fatalf("stock left at %d", products.items["coat"].Stock)
	}
	if len(orders.items) != 0 {
		t.Fatal("order was stored after failure")
	}

	orders.fail = false
	order, err := service.Place(context.Background(), "north", PlaceOrderInput{
		Customer: models.Customer{Name: "Ada Lovelace", Email: "ada@example.com", Address: "12 Analytical", City: "London"},
		Items:    []ItemInput{{ProductID: "coat", Quantity: 1, UnitPriceCents: 2500}},
	})
	if err != nil {
		t.Fatal(err)
	}
	products.items["coat"] = func() models.Product {
		p := products.items["coat"]
		p.PriceCents = 100
		return p
	}()
	if order.Items[0].UnitPriceCents != 2500 || order.TotalCents != 2500 {
		t.Fatalf("snapshot = %+v", order.Items)
	}
	if products.items["coat"].Stock != 1 {
		t.Fatalf("stock = %d", products.items["coat"].Stock)
	}

	cancelled, err := service.UpdateStatus(context.Background(), "user-1", order.ID, models.OrderCancelled)
	if err != nil {
		t.Fatal(err)
	}
	if cancelled.Status != models.OrderCancelled || products.items["coat"].Stock != 2 {
		t.Fatalf("status %s stock %d", cancelled.Status, products.items["coat"].Stock)
	}
	if _, err := service.UpdateStatus(context.Background(), "user-1", order.ID, models.OrderShipped); err == nil {
		t.Fatal("cancelled order changed status")
	}
	if products.items["coat"].Stock != 2 {
		t.Fatalf("stock changed after rejected update: %d", products.items["coat"].Stock)
	}
}

func TestPlaceOrderChargesTheChosenVariant(t *testing.T) {
	stores := &memStores{items: []models.Store{{ID: "store-1", OwnerID: "user-1", Slug: "north", Published: true}}}
	products := &memProducts{items: map[string]models.Product{
		"coat": {
			ID: "coat", StoreID: "store-1", Name: "Coat", PriceCents: 2000, Stock: 4, Published: true,
			OptionNames: []string{"Color"},
			Variants: []models.Variant{
				{ID: "black", Options: []models.VariantOption{{Name: "Color", Value: "Black"}}, SKU: "C-B", PriceCents: 2000, Stock: 3},
				{ID: "sand", Options: []models.VariantOption{{Name: "Color", Value: "Sand"}}, SKU: "C-S", PriceCents: 2200, Stock: 1},
			},
		},
	}}
	orders := &memOrders{}
	service := NewOrderService(stores, products, orders)
	order, err := service.Place(context.Background(), "north", PlaceOrderInput{
		Customer: models.Customer{Name: "Ada Lovelace", Email: "ada@example.com", Address: "12 Analytical", City: "London"},
		Items:    []ItemInput{{ProductID: "coat", VariantID: "sand", Quantity: 1, UnitPriceCents: 2200}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.Items[0].OptionLabel != "Sand" || order.Items[0].SKU != "C-S" || order.TotalCents != 2200 {
		t.Fatalf("snapshot = %+v", order.Items)
	}
	coat := products.items["coat"]
	if coat.Variants[0].Stock != 3 || coat.Variants[1].Stock != 0 || coat.Stock != 3 {
		t.Fatalf("stock = %+v total %d", coat.Variants, coat.Stock)
	}
	if _, err := service.Place(context.Background(), "north", PlaceOrderInput{
		Customer: models.Customer{Name: "Ada Lovelace", Email: "ada@example.com", Address: "12 Analytical", City: "London"},
		Items:    []ItemInput{{ProductID: "coat", Quantity: 1, UnitPriceCents: 2000}},
	}); err == nil {
		t.Fatal("expected a version choice")
	}
	if products.items["coat"].Stock != 3 {
		t.Fatalf("stock changed without a version: %d", products.items["coat"].Stock)
	}
}

func TestPlaceOrderRejectsInsufficientStock(t *testing.T) {
	stores := &memStores{items: []models.Store{{ID: "store-1", Slug: "north", Published: true}}}
	products := &memProducts{items: map[string]models.Product{
		"coat": {ID: "coat", StoreID: "store-1", Name: "Coat", PriceCents: 2500, Stock: 1, Published: true},
	}}
	orders := &memOrders{}
	service := NewOrderService(stores, products, orders)
	_, err := service.Place(context.Background(), "north", PlaceOrderInput{
		Customer: models.Customer{Name: "Ada Lovelace", Email: "ada@example.com", Address: "12 Analytical", City: "London"},
		Items:    []ItemInput{{ProductID: "coat", Quantity: 4, UnitPriceCents: 2500}},
	})
	if err == nil {
		t.Fatal("expected stock error")
	}
	if products.items["coat"].Stock != 1 || len(orders.items) != 0 {
		t.Fatalf("stock %d orders %d", products.items["coat"].Stock, len(orders.items))
	}
}

type memStores struct{ items []models.Store }

func (m *memStores) Create(context.Context, models.Store) (models.Store, error) {
	return models.Store{}, errors.New("not implemented")
}
func (m *memStores) FindByOwner(_ context.Context, ownerID string) (models.Store, error) {
	for _, store := range m.items {
		if store.OwnerID == ownerID {
			return store, nil
		}
	}
	return models.Store{}, interfaces.ErrNotFound
}
func (m *memStores) FindBySlug(_ context.Context, slug string) (models.Store, error) {
	for _, store := range m.items {
		if store.Slug == slug {
			return store, nil
		}
	}
	return models.Store{}, interfaces.ErrNotFound
}
func (m *memStores) UpdateProfile(context.Context, models.Store) error { return nil }
func (m *memStores) SetPublished(context.Context, string, string, bool, time.Time) error {
	return nil
}
func (m *memStores) UpdateDraft(context.Context, string, string, models.StorefrontConfig, *models.DesignNote, time.Time) error {
	return nil
}
func (m *memStores) PublishDraft(context.Context, string, string, time.Time) error { return nil }

type memProducts struct{ items map[string]models.Product }

func (m *memProducts) Create(context.Context, models.Product) (models.Product, error) {
	return models.Product{}, errors.New("not implemented")
}
func (m *memProducts) Update(context.Context, models.Product) (models.Product, error) {
	return models.Product{}, errors.New("not implemented")
}
func (m *memProducts) Delete(context.Context, string, string) error { return nil }
func (m *memProducts) FindByID(_ context.Context, storeID, id string) (models.Product, error) {
	product, ok := m.items[id]
	if !ok || product.StoreID != storeID {
		return models.Product{}, interfaces.ErrNotFound
	}
	return product, nil
}
func (m *memProducts) List(context.Context, string, interfaces.ProductFilter) ([]models.Product, error) {
	return nil, nil
}
func (m *memProducts) AddStock(_ context.Context, storeID, id, variantID string, delta int) error {
	product, ok := m.items[id]
	if !ok || product.StoreID != storeID {
		return interfaces.ErrNotFound
	}
	if variantID == "" {
		variantID = id
	}
	if len(product.Variants) == 0 {
		if variantID != id && variantID != product.ID {
			return interfaces.ErrNotFound
		}
		if product.Stock+delta < 0 {
			return interfaces.ErrInsufficientStock
		}
		product.Stock += delta
		m.items[id] = product
		return nil
	}
	found := false
	for i := range product.Variants {
		if product.Variants[i].ID != variantID {
			continue
		}
		found = true
		if product.Variants[i].Stock+delta < 0 {
			return interfaces.ErrInsufficientStock
		}
		product.Variants[i].Stock += delta
	}
	if !found {
		return interfaces.ErrNotFound
	}
	total := 0
	for _, variant := range product.Variants {
		total += variant.Stock
	}
	product.Stock = total
	m.items[id] = product
	return nil
}
func (m *memProducts) SetStock(context.Context, string, string, string, int, time.Time) error {
	return nil
}
func (m *memProducts) SKUInUse(context.Context, string, string, string) (bool, error) {
	return false, nil
}
func (m *memProducts) Totals(context.Context, string) (int, int, error) { return 0, 0, nil }

type memOrders struct {
	items []models.Order
	fail  bool
}

func (m *memOrders) Create(_ context.Context, order models.Order) (models.Order, error) {
	if m.fail {
		return models.Order{}, errors.New("db down")
	}
	order.ID = "order-1"
	m.items = append(m.items, order)
	return order, nil
}
func (m *memOrders) FindByID(_ context.Context, storeID, id string) (models.Order, error) {
	for _, order := range m.items {
		if order.ID == id && order.StoreID == storeID {
			return order, nil
		}
	}
	return models.Order{}, interfaces.ErrNotFound
}
func (m *memOrders) List(context.Context, string, interfaces.OrderFilter) ([]models.Order, error) {
	return nil, nil
}
func (m *memOrders) UpdateStatus(_ context.Context, storeID, id string, status models.OrderStatus, at time.Time) (models.Order, error) {
	for i, order := range m.items {
		if order.ID == id && order.StoreID == storeID {
			if order.Status == models.OrderCancelled {
				return models.Order{}, interfaces.ErrConflict
			}
			order.Status = status
			order.UpdatedAt = at
			m.items[i] = order
			return order, nil
		}
	}
	return models.Order{}, interfaces.ErrNotFound
}
func (m *memOrders) Summary(context.Context, string, time.Time, models.OrderStatus) (int64, int, error) {
	return 0, 0, nil
}
func (m *memOrders) RevenueByDay(context.Context, string, time.Time, models.OrderStatus) ([]models.RevenuePoint, error) {
	return nil, nil
}
func (m *memOrders) BestSellers(context.Context, string, time.Time, models.OrderStatus, int) ([]models.BestSeller, error) {
	return nil, nil
}
