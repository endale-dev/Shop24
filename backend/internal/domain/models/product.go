package models

import "time"

type Product struct {
	ID          string    `json:"id"`
	StoreID     string    `json:"storeId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceCents  int64     `json:"priceCents"`
	Images      []string  `json:"images"`
	Category    string    `json:"category"`
	SKU         string    `json:"sku"`
	Stock       int       `json:"stock"`
	LowStock    bool      `json:"lowStock"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
