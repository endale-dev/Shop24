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

type OrderService struct {
	stores   interfaces.StoreRepository
	products interfaces.ProductRepository
	orders   interfaces.OrderRepository
}

func NewOrderService(stores interfaces.StoreRepository, products interfaces.ProductRepository, orders interfaces.OrderRepository) *OrderService {
	return &OrderService{stores: stores, products: products, orders: orders}
}

type ItemInput struct {
	ProductID      string
	Quantity       int
	UnitPriceCents int64
}

type PlaceOrderInput struct {
	Customer models.Customer
	Items    []ItemInput
}

func (s *OrderService) Place(ctx context.Context, slug string, in PlaceOrderInput) (models.Order, error) {
	store, err := publishedStore(ctx, s.stores, slug)
	if err != nil {
		return models.Order{}, err
	}
	customer, err := validateCustomer(in.Customer)
	if err != nil {
		return models.Order{}, err
	}
	merged := mergeItems(in.Items)
	if len(merged) == 0 || len(merged) > 20 {
		return models.Order{}, apperr.New(http.StatusBadRequest, "invalid_items", "Add between 1 and 20 products")
	}
	items := make([]models.OrderItem, 0, len(merged))
	var subtotal int64
	for _, item := range merged {
		if item.Quantity < 1 || item.Quantity > 99 {
			return models.Order{}, apperr.New(http.StatusBadRequest, "invalid_quantity", "Quantity must be between 1 and 99")
		}
		product, err := s.products.FindByID(ctx, store.ID, item.ProductID)
		if errors.Is(err, interfaces.ErrNotFound) || (err == nil && !product.Published) {
			return models.Order{}, apperr.New(http.StatusConflict, "product_unavailable", "A product in this order is no longer available")
		}
		if err != nil {
			return models.Order{}, err
		}
		if item.UnitPriceCents != product.PriceCents {
			return models.Order{}, apperr.New(http.StatusConflict, "price_changed", "The price of "+product.Name+" changed. Review the cart and try again.")
		}
		if product.Stock < item.Quantity {
			return models.Order{}, apperr.New(http.StatusConflict, "insufficient_stock", "Not enough stock for "+product.Name)
		}
		line := product.PriceCents * int64(item.Quantity)
		subtotal += line
		items = append(items, models.OrderItem{
			ProductID: product.ID, Name: product.Name, SKU: product.SKU,
			UnitPriceCents: product.PriceCents, Quantity: item.Quantity, LineTotalCents: line,
		})
	}

	held := make([]ItemInput, 0, len(items))
	for _, item := range items {
		if err := s.products.AddStock(ctx, store.ID, item.ProductID, -item.Quantity); err != nil {
			s.release(ctx, store.ID, held)
			if errors.Is(err, interfaces.ErrInsufficientStock) {
				return models.Order{}, apperr.New(http.StatusConflict, "insufficient_stock", "Not enough stock for "+item.Name)
			}
			if errors.Is(err, interfaces.ErrNotFound) {
				return models.Order{}, apperr.New(http.StatusConflict, "product_unavailable", "A product in this order is no longer available")
			}
			return models.Order{}, err
		}
		held = append(held, ItemInput{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	now := time.Now().UTC()
	created, err := s.orders.Create(ctx, models.Order{
		StoreID: store.ID, Status: models.OrderPending, Customer: customer, Items: items,
		SubtotalCents: subtotal, TotalCents: subtotal, PaymentMethod: "cod",
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		s.release(ctx, store.ID, held)
		return models.Order{}, err
	}
	return created, nil
}

func (s *OrderService) List(ctx context.Context, ownerID, status string) ([]models.Order, error) {
	if status != "" && !models.ValidOrderStatus(models.OrderStatus(status)) {
		return nil, apperr.New(http.StatusBadRequest, "invalid_status", "Choose a valid order status")
	}
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return nil, err
	}
	orders, err := s.orders.List(ctx, store.ID, interfaces.OrderFilter{Status: status, Limit: 100})
	if err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []models.Order{}
	}
	return orders, nil
}

func (s *OrderService) Get(ctx context.Context, ownerID, orderID string) (models.Order, error) {
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Order{}, err
	}
	order, err := s.orders.FindByID(ctx, store.ID, orderID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Order{}, apperr.New(http.StatusNotFound, "not_found", "Order not found")
	}
	return order, err
}

func (s *OrderService) UpdateStatus(ctx context.Context, ownerID, orderID string, status models.OrderStatus) (models.Order, error) {
	if !models.ValidOrderStatus(status) {
		return models.Order{}, apperr.New(http.StatusBadRequest, "invalid_status", "Choose a valid order status")
	}
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Order{}, err
	}
	current, err := s.orders.FindByID(ctx, store.ID, orderID)
	if errors.Is(err, interfaces.ErrNotFound) {
		return models.Order{}, apperr.New(http.StatusNotFound, "not_found", "Order not found")
	}
	if err != nil {
		return models.Order{}, err
	}
	if current.Status == models.OrderCancelled {
		return models.Order{}, apperr.New(http.StatusConflict, "order_closed", "Cancelled orders cannot change status")
	}
	if current.Status == status {
		return current, nil
	}
	if status == models.OrderCancelled {
		if err := s.restore(ctx, store.ID, current.Items); err != nil {
			return models.Order{}, err
		}
	}
	updated, err := s.orders.UpdateStatus(ctx, store.ID, orderID, status, time.Now().UTC())
	if err != nil && status == models.OrderCancelled {
		_ = s.reserveAgain(ctx, store.ID, current.Items)
	}
	if errors.Is(err, interfaces.ErrConflict) {
		return models.Order{}, apperr.New(http.StatusConflict, "order_closed", "Cancelled orders cannot change status")
	}
	return updated, err
}

func (s *OrderService) restore(ctx context.Context, storeID string, items []models.OrderItem) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	restored := make([]models.OrderItem, 0, len(items))
	for _, item := range items {
		err := s.products.AddStock(ctx, storeID, item.ProductID, item.Quantity)
		if errors.Is(err, interfaces.ErrNotFound) {
			continue
		}
		if err != nil {
			_ = s.reserveAgain(ctx, storeID, restored)
			return err
		}
		restored = append(restored, item)
	}
	return nil
}

func (s *OrderService) reserveAgain(ctx context.Context, storeID string, items []models.OrderItem) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	for i := len(items) - 1; i >= 0; i-- {
		if err := s.products.AddStock(ctx, storeID, items[i].ProductID, -items[i].Quantity); err != nil && !errors.Is(err, interfaces.ErrNotFound) {
			return err
		}
	}
	return nil
}

func (s *OrderService) release(ctx context.Context, storeID string, held []ItemInput) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	for i := len(held) - 1; i >= 0; i-- {
		_ = s.products.AddStock(ctx, storeID, held[i].ProductID, held[i].Quantity)
	}
}

func mergeItems(items []ItemInput) []ItemInput {
	order := make([]string, 0, len(items))
	qty := map[string]int{}
	price := map[string]int64{}
	for _, item := range items {
		if item.ProductID == "" {
			continue
		}
		if _, ok := qty[item.ProductID]; !ok {
			order = append(order, item.ProductID)
			price[item.ProductID] = item.UnitPriceCents
		}
		qty[item.ProductID] += item.Quantity
	}
	out := make([]ItemInput, 0, len(order))
	for _, id := range order {
		out = append(out, ItemInput{ProductID: id, Quantity: qty[id], UnitPriceCents: price[id]})
	}
	return out
}

func validateCustomer(customer models.Customer) (models.Customer, error) {
	name, err := input.Text(customer.Name, "Name", 2, 80)
	if err != nil {
		return models.Customer{}, apperr.New(http.StatusBadRequest, "invalid_customer", err.Error())
	}
	email, err := input.Email(customer.Email)
	if err != nil {
		return models.Customer{}, apperr.New(http.StatusBadRequest, "invalid_customer", err.Error())
	}
	phone, err := input.Optional(customer.Phone, 30)
	if err != nil {
		return models.Customer{}, apperr.New(http.StatusBadRequest, "invalid_customer", "Phone number is too long")
	}
	address, err := input.Text(customer.Address, "Address", 4, 160)
	if err != nil {
		return models.Customer{}, apperr.New(http.StatusBadRequest, "invalid_customer", err.Error())
	}
	city, err := input.Text(customer.City, "City", 2, 60)
	if err != nil {
		return models.Customer{}, apperr.New(http.StatusBadRequest, "invalid_customer", err.Error())
	}
	notes, err := input.Optional(customer.Notes, 300)
	if err != nil {
		return models.Customer{}, apperr.New(http.StatusBadRequest, "invalid_customer", "Notes must be 300 characters or less")
	}
	return models.Customer{Name: name, Email: email, Phone: phone, Address: address, City: city, Notes: notes}, nil
}
