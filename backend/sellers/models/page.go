package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PageBlock represents a single block in the page tree.
type PageBlock struct {
	ID           string                 `bson:"id" json:"id"`
	ComponentKey string                 `bson:"component_key" json:"componentKey"`
	Props        map[string]interface{} `bson:"props,omitempty" json:"props,omitempty"`
	Children     []PageBlock            `bson:"children,omitempty" json:"children,omitempty"`
}

// Page represents a storefront page created by sellers.
type Page struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	StoreID    primitive.ObjectID     `bson:"store_id" json:"storeId"`
	Path       string                 `bson:"path" json:"path"` // e.g. "/", "/about"
	Title      string                 `bson:"title" json:"title"`
	SEO        map[string]interface{} `bson:"seo,omitempty" json:"seo,omitempty"`
	Visibility string                 `bson:"visibility" json:"visibility"` // public|authenticated|roles
	Status     string                 `bson:"status" json:"status"`         // draft|published
	Blocks     []PageBlock            `bson:"blocks" json:"blocks"`
	Snapshot   map[string]interface{} `bson:"snapshot,omitempty" json:"snapshot,omitempty"` // published snapshot
	CreatedAt  time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updatedAt"`
}
