package models

import (
	"strings"
	"time"
)

type VariantOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Variant struct {
	ID         string          `json:"id"`
	Options    []VariantOption `json:"options"`
	SKU        string          `json:"sku"`
	PriceCents int64           `json:"priceCents"`
	Stock      int             `json:"stock"`
	ImageURL   string          `json:"imageUrl"`
}

func (v Variant) Label() string {
	parts := make([]string, 0, len(v.Options))
	for _, option := range v.Options {
		if option.Value != "" {
			parts = append(parts, option.Value)
		}
	}
	return strings.Join(parts, " / ")
}

type Product struct {
	ID          string    `json:"id"`
	StoreID     string    `json:"storeId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceCents  int64     `json:"priceCents"`
	Images      []string  `json:"images"`
	Category    string    `json:"category"`
	OptionNames []string  `json:"optionNames"`
	Variants    []Variant `json:"variants"`
	SKU         string    `json:"sku"`
	Stock       int       `json:"stock"`
	LowStock    bool      `json:"lowStock"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (p Product) WithVariants() Product {
	if p.OptionNames == nil {
		p.OptionNames = []string{}
	}
	if len(p.Variants) == 0 {
		p.OptionNames = []string{}
		p.Variants = []Variant{{
			ID: p.ID, SKU: p.SKU, PriceCents: p.PriceCents, Stock: p.Stock, Options: []VariantOption{},
		}}
	}
	for i := range p.Variants {
		if p.Variants[i].Options == nil {
			p.Variants[i].Options = []VariantOption{}
		}
	}
	minPrice := p.Variants[0].PriceCents
	total := 0
	for _, variant := range p.Variants {
		if variant.PriceCents < minPrice {
			minPrice = variant.PriceCents
		}
		total += variant.Stock
	}
	p.PriceCents = minPrice
	p.Stock = total
	p.SKU = p.Variants[0].SKU
	return p
}
