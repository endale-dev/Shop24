package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
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
	VariantID      string
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
		variant, err := pickVariant(product, item.VariantID)
		if err != nil {
			return models.Order{}, err
		}
		label := variantDisplay(product.Name, variant)
		if item.UnitPriceCents != variant.PriceCents {
			return models.Order{}, apperr.New(http.StatusConflict, "price_changed", "The price of "+label+" changed. Review the cart and try again.")
		}
		if variant.Stock < item.Quantity {
			return models.Order{}, apperr.New(http.StatusConflict, "insufficient_stock", "Not enough stock for "+label)
		}
		line := variant.PriceCents * int64(item.Quantity)
		subtotal += line
		items = append(items, models.OrderItem{
			ProductID: product.ID, VariantID: variant.ID, Name: product.Name, OptionLabel: variant.Label(), SKU: variant.SKU,
			UnitPriceCents: variant.PriceCents, Quantity: item.Quantity, LineTotalCents: line,
		})
	}

	held := make([]ItemInput, 0, len(items))
	for _, item := range items {
		if err := s.products.AddStock(ctx, store.ID, item.ProductID, item.VariantID, -item.Quantity); err != nil {
			relErr := s.release(ctx, store.ID, held)
			if errors.Is(err, interfaces.ErrInsufficientStock) {
				err = apperr.New(http.StatusConflict, "insufficient_stock", "Not enough stock for "+item.Name)
			} else if errors.Is(err, interfaces.ErrNotFound) {
				err = apperr.New(http.StatusConflict, "product_unavailable", "A product in this order is no longer available")
			}
			if relErr != nil {
				err = errors.Join(err, relErr)
			}
			return models.Order{}, err
		}
		held = append(held, ItemInput{ProductID: item.ProductID, VariantID: item.VariantID, Quantity: item.Quantity})
	}

	now := time.Now().UTC()
	created, err := s.orders.Create(ctx, models.Order{
		StoreID: store.ID, Status: models.OrderPending, Customer: customer, Items: items,
		SubtotalCents: subtotal, TotalCents: subtotal, PaymentMethod: "cod",
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		if relErr := s.release(ctx, store.ID, held); relErr != nil {
			err = errors.Join(err, relErr)
		}
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
		if undoErr := s.reserveAgain(ctx, store.ID, current.Items); undoErr != nil {
			err = errors.Join(err, undoErr)
		}
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
		err := s.products.AddStock(ctx, storeID, item.ProductID, item.VariantID, item.Quantity)
		if errors.Is(err, interfaces.ErrNotFound) {
			continue
		}
		if err != nil {
			if undoErr := s.reserveAgain(ctx, storeID, restored); undoErr != nil {
				return errors.Join(err, undoErr)
			}
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
		if err := s.products.AddStock(ctx, storeID, items[i].ProductID, items[i].VariantID, -items[i].Quantity); err != nil && !errors.Is(err, interfaces.ErrNotFound) {
			return err
		}
	}
	return nil
}

func (s *OrderService) release(ctx context.Context, storeID string, held []ItemInput) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	var first error
	for i := len(held) - 1; i >= 0; i-- {
		if err := s.products.AddStock(ctx, storeID, held[i].ProductID, held[i].VariantID, held[i].Quantity); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func mergeItems(items []ItemInput) []ItemInput {
	order := make([]string, 0, len(items))
	qty := map[string]int{}
	price := map[string]int64{}
	variant := map[string]string{}
	for _, item := range items {
		if item.ProductID == "" {
			continue
		}
		key := item.ProductID + "\x00" + item.VariantID
		if _, ok := qty[key]; !ok {
			order = append(order, key)
			price[key] = item.UnitPriceCents
			variant[key] = item.VariantID
		}
		qty[key] += item.Quantity
	}
	out := make([]ItemInput, 0, len(order))
	for _, key := range order {
		productID, _, _ := strings.Cut(key, "\x00")
		out = append(out, ItemInput{ProductID: productID, VariantID: variant[key], Quantity: qty[key], UnitPriceCents: price[key]})
	}
	return out
}

func pickVariant(product models.Product, variantID string) (models.Variant, error) {
	product = product.WithVariants()
	if variantID == "" {
		if len(product.Variants) == 1 {
			return product.Variants[0], nil
		}
		return models.Variant{}, apperr.New(http.StatusBadRequest, "choose_variant", "Choose a version of "+product.Name)
	}
	for _, variant := range product.Variants {
		if variant.ID == variantID {
			return variant, nil
		}
	}
	return models.Variant{}, apperr.New(http.StatusConflict, "product_unavailable", "A product in this order is no longer available")
}

func variantDisplay(name string, variant models.Variant) string {
	if label := variant.Label(); label != "" {
		return name + " (" + label + ")"
	}
	return name
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
