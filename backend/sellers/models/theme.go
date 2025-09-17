package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Theme represents design tokens and default component props for a store.
type Theme struct {
	ID                primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	StoreID           primitive.ObjectID     `bson:"store_id" json:"storeId"`
	Name              string                 `bson:"name" json:"name"`
	Description       string                 `bson:"description,omitempty" json:"description,omitempty"`
	Tokens            map[string]interface{} `bson:"tokens" json:"tokens"` // colors, typography, spacing
	ComponentDefaults map[string]interface{} `bson:"component_defaults,omitempty" json:"componentDefaults,omitempty"`
	IsActive          bool                   `bson:"is_active" json:"isActive"`
	CreatedAt         time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time              `bson:"updated_at" json:"updatedAt"`
}
