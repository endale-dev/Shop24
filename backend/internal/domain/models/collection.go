package models

import "time"

type Collection struct {
	ID          string    `json:"id"`
	StoreID     string    `json:"storeId"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	ProductIDs  []string  `json:"productIds"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
