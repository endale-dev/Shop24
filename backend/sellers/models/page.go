package models

import (
	"time"

	sharedModels "github.com/Endale2/DRPS/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PageBlock alias for backward compatibility
type PageBlock = sharedModels.Block

// Page represents a storefront page created by sellers.
type Page struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	StoreID    primitive.ObjectID     `bson:"store_id" json:"storeId"`
	Path       string                 `bson:"path" json:"path"` // e.g. "/", "/about"
	Title      string                 `bson:"title" json:"title"`
	SEO        map[string]interface{} `bson:"seo,omitempty" json:"seo,omitempty"`
	Visibility string                 `bson:"visibility" json:"visibility"`                 // public|authenticated|roles
	Status     string                 `bson:"status" json:"status"`                         // draft|published
	Blocks     []sharedModels.Block   `bson:"blocks" json:"blocks"`                         // Using shared.Block
	Snapshot   map[string]interface{} `bson:"snapshot,omitempty" json:"snapshot,omitempty"` // published snapshot
	CreatedAt  time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updatedAt"`
}
