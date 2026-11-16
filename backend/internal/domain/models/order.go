package models

import "time"

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
)

func ValidOrderStatus(status OrderStatus) bool {
	switch status {
	case OrderPending, OrderConfirmed, OrderProcessing, OrderShipped, OrderDelivered, OrderCancelled:
		return true
	default:
		return false
	}
}

type OrderItem struct {
	ProductID      string `json:"productId"`
	Name           string `json:"name"`
	SKU            string `json:"sku"`
	UnitPriceCents int64  `json:"unitPriceCents"`
	Quantity       int    `json:"quantity"`
	LineTotalCents int64  `json:"lineTotalCents"`
}

type Customer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	City    string `json:"city"`
	Notes   string `json:"notes"`
}

type Order struct {
	ID            string      `json:"id"`
	StoreID       string      `json:"storeId"`
	Status        OrderStatus `json:"status"`
	Customer      Customer    `json:"customer"`
	Items         []OrderItem `json:"items"`
	SubtotalCents int64       `json:"subtotalCents"`
	TotalCents    int64       `json:"totalCents"`
	PaymentMethod string      `json:"paymentMethod"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}
