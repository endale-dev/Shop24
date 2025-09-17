// backend/models/theme.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Theme represents a design theme (tokens + defaults) for a store
type Theme struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ShopID primitive.ObjectID `bson:"shop_id,omitempty"         json:"shop_id"` // tenant isolation
	Name        string                 `bson:"name" json:"name"`
	Version     string                 `bson:"version" json:"version"`
	Description string                 `bson:"description,omitempty" json:"description,omitempty"`

	// Design tokens (CSS variables style)
	Tokens map[string]interface{} `bson:"tokens" json:"tokens"`

	// Default props per componentKey
	ComponentsDefaults map[string]map[string]interface{} `bson:"componentsDefaults" json:"componentsDefaults"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy" json:"updatedBy"`
}
