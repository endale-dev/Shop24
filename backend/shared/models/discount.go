package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DiscountCategory is *where* it applies:
type DiscountCategory string

const (
	DiscountCategoryProduct DiscountCategory = "product" // specific products or variants
	// TODO: Reintroduce these categories in future versions
	// DiscountCategoryOrder    DiscountCategory = "order"       // entire cart
	// DiscountCategoryShipping DiscountCategory = "shipping"    // free shipping, etc.
	// DiscountCategoryBuyXGetY DiscountCategory = "buy_x_get_y" // buy X get Y
)

// DiscountType is *how much*:
type DiscountType string

const (
	DiscountTypeFixed      DiscountType = "fixed"      // $X off
	DiscountTypePercentage DiscountType = "percentage" // X% off
)

// DiscountEligibilityType defines how customer eligibility is determined
type DiscountEligibilityType string

const (
	DiscountEligibilityAll      DiscountEligibilityType = "all"      // Everyone can use
	DiscountEligibilitySpecific DiscountEligibilityType = "specific" // Only specific customers
	DiscountEligibilitySegment  DiscountEligibilityType = "segment"  // Only customers in specific segments
)

// DiscountUsage tracks usage of a discount
type DiscountUsage struct {
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	UsageCount int                `bson:"usage_count" json:"usage_count"`
	LastUsedAt time.Time          `bson:"last_used_at" json:"last_used_at"`
	TotalSpent float64            `bson:"total_spent" json:"total_spent"`
}

// Discount represents a discount offer in the system.
// NOTE: All business logic for discounts has been moved to shared/services/discount_service.go
// This model is now a pure data structure.
type Discount struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`

	// *where* it applies - only product-level discounts
	Category DiscountCategory `bson:"category" json:"category"`
	// *how much* off
	Type  DiscountType `bson:"type,omitempty" json:"type,omitempty"`
	Value float64      `bson:"value,omitempty" json:"value,omitempty"`

	ShopID   primitive.ObjectID `bson:"shop_id" json:"shop_id"`
	SellerID primitive.ObjectID `bson:"seller_id" json:"seller_id"`

	// Product-level targets only
	AppliesToProducts []primitive.ObjectID `bson:"applies_to_products,omitempty" json:"applies_to_products,omitempty"`
	AppliesToVariants []primitive.ObjectID `bson:"applies_to_variants,omitempty" json:"applies_to_variants,omitempty"`
	// TODO: Reintroduce collection-level discounts in future versions
	// AppliesToCollections []primitive.ObjectID `bson:"applies_to_collections,omitempty" json:"applies_to_collections,omitempty"`

	// Enhanced customer eligibility
	EligibilityType    DiscountEligibilityType `bson:"eligibility_type" json:"eligibility_type"`
	AllowedCustomerIDs []primitive.ObjectID    `bson:"allowed_customers,omitempty" json:"allowed_customers,omitempty"`
	AllowedSegmentIDs  []primitive.ObjectID    `bson:"allowed_segments,omitempty" json:"allowed_segments,omitempty"`

	// Usage limits
	UsageLimit       *int            `bson:"usage_limit,omitempty" json:"usage_limit,omitempty"`
	PerCustomerLimit *int            `bson:"per_customer_limit,omitempty" json:"per_customer_limit,omitempty"`
	CurrentUsage     int             `bson:"current_usage" json:"current_usage"`
	UsageTracking    []DiscountUsage `bson:"usage_tracking,omitempty" json:"usage_tracking,omitempty"`

	// TODO: Reintroduce Buy X Get Y functionality in future versions
	// BuyProductIDs []primitive.ObjectID `bson:"buy_product_ids,omitempty" json:"buy_product_ids,omitempty"`
	// BuyQuantity   *int                 `bson:"buy_quantity,omitempty" json:"buy_quantity,omitempty"`
	// GetProductIDs []primitive.ObjectID `bson:"get_product_ids,omitempty" json:"get_product_ids,omitempty"`
	// GetQuantity   *int                 `bson:"get_quantity,omitempty" json:"get_quantity,omitempty"`

	// TODO: Reintroduce order-level functionality in future versions
	// MinimumOrderSubtotal *float64 `bson:"minimum_order_subtotal,omitempty" json:"minimum_order_subtotal,omitempty"`

	// TODO: Reintroduce shipping-level functionality in future versions
	// FreeShipping                bool     `bson:"free_shipping,omitempty" json:"free_shipping,omitempty"`
	// MinimumOrderForFreeShipping *float64 `bson:"minimum_free_shipping,omitempty" json:"minimum_free_shipping,omitempty"`

	StartAt time.Time `bson:"start_at" json:"start_at"`
	EndAt   time.Time `bson:"end_at" json:"end_at"`
	Active  bool      `bson:"active" json:"active"`

	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
