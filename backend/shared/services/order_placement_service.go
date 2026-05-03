// shared/services/order_placement_service.go
package services

import (
	"errors"
	"time"

	"github.com/Endale2/DRPS/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItemRequest represents the minimal data needed from client for order placement
type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id,omitempty"` // optional, empty string for no variant
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// PlaceOrderRequest represents the complete request structure for placing an order
type PlaceOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
}

// OrderPlacementResult contains the created order and additional details
type OrderPlacementResult struct {
	Order               *models.Order
	ItemDiscountDetails []map[string]interface{}
}

// Errors for order placement
var (
	ErrInvalidProductID  = errors.New("invalid product ID")
	ErrInvalidVariantID  = errors.New("invalid variant ID")
	ErrProductNotFound   = errors.New("product not found")
	ErrProductNotInShop  = errors.New("product does not belong to this shop")
	ErrVariantNotFound   = errors.New("variant not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
)

// PlaceOrderFromRequest handles the complete order placement flow with server-side validation and pricing
func PlaceOrderFromRequest(shopID, customerID primitive.ObjectID, req PlaceOrderRequest) (*OrderPlacementResult, error) {
	// Get customer segments for discount validation
	customerSegmentIDs, err := GetCustomerSegmentIDs(shopID, customerID)
	if err != nil {
		customerSegmentIDs = []primitive.ObjectID{}
	}

	// Initialize order variables
	var orderItems []models.OrderItem
	var itemDiscountDetails []map[string]interface{}
	subtotal := 0.0
	totalDiscount := 0.0
	appliedDiscountIDsMap := make(map[primitive.ObjectID]struct{})

	// Process each item - ALL pricing calculated server-side
	for _, itemReq := range req.Items {
		// Validate product ID
		productID, err := primitive.ObjectIDFromHex(itemReq.ProductID)
		if err != nil {
			return nil, ErrInvalidProductID
		}

		// Validate variant ID if provided
		variantID := primitive.NilObjectID
		if itemReq.VariantID != "" {
			variantID, err = primitive.ObjectIDFromHex(itemReq.VariantID)
			if err != nil {
				return nil, ErrInvalidVariantID
			}
		}

		// Validate quantity
		if itemReq.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}

		// Fetch product from database - server-side validation
		product, err := GetProductByIDService(itemReq.ProductID)
		if err != nil || product == nil {
			return nil, ErrProductNotFound
		}

		// Validate product belongs to this shop
		if product.ShopID != shopID {
			return nil, ErrProductNotInShop
		}

		// Determine unit price and product details - server-side calculation
		var unitPrice float64
		var productName string
		var productImage string
		var stock int
		var orderVariantID primitive.ObjectID

		// Check if product has variants
		hasRealVariants := len(product.Variants) > 0

		if hasRealVariants && !variantID.IsZero() {
			// Variant-specific pricing
			var foundVariant *models.Variant
			for i := range product.Variants {
				if product.Variants[i].VariantID == variantID {
					foundVariant = &product.Variants[i]
					break
				}
			}
			if foundVariant == nil {
				return nil, ErrVariantNotFound
			}
			unitPrice = foundVariant.Price
			productName = product.Name
			if len(foundVariant.Options) > 0 {
				productName += " - "
				for i, opt := range foundVariant.Options {
					if i > 0 {
						productName += ", "
					}
					productName += opt.Name + ": " + opt.Value
				}
			}
			if foundVariant.Image != "" {
				productImage = foundVariant.Image
			} else {
				productImage = product.MainImage
			}
			stock = foundVariant.Stock
			orderVariantID = variantID
		} else {
			// Product-level pricing
			unitPrice = product.Price
			productName = product.Name
			productImage = product.MainImage
			stock = product.Stock
			orderVariantID = primitive.NilObjectID
		}

		// Validate stock availability
		if stock < itemReq.Quantity {
			return nil, ErrInsufficientStock
		}

		// Calculate line total server-side
		lineTotal := unitPrice * float64(itemReq.Quantity)
		subtotal += lineTotal

		// Get collection IDs for discount eligibility
		collectionIDs, err := GetCollectionIDsForProduct(productID)
		if err != nil {
			collectionIDs = []primitive.ObjectID{}
		}

		// Get active discounts for this product/variant - ALL discount logic server-side
		discounts, err := GetActiveDiscountsForProductService(shopID, productID, variantID, collectionIDs)
		if err != nil {
			discounts = []models.Discount{}
		}

		// Find best eligible discount - server-side calculation only
		bestDiscount, _ := GetBestEligibleDiscountForProduct(productID, variantID, customerID, customerSegmentIDs, discounts)

		var bestSavings float64
		if bestDiscount != nil {
			// Calculate actual savings for this specific order
			bestSavings = CalculateDiscountAmountForQuantity(bestDiscount, unitPrice, itemReq.Quantity)
		}

		// Apply discount if found - server-side calculation
		itemDiscountAmount := 0.0
		appliedDiscountIDs := []primitive.ObjectID{}
		if bestDiscount != nil && bestSavings > 0 {
			itemDiscountAmount = bestSavings
			appliedDiscountIDs = append(appliedDiscountIDs, bestDiscount.ID)
			appliedDiscountIDsMap[bestDiscount.ID] = struct{}{}
		}

		totalDiscount += itemDiscountAmount
		finalLineTotal := lineTotal - itemDiscountAmount
		if finalLineTotal < 0 {
			finalLineTotal = 0
		}

		// Store discount details for response
		itemDiscountDetails = append(itemDiscountDetails, map[string]interface{}{
			"product_id":           productID.Hex(),
			"variant_id":           variantID.Hex(),
			"unit_price":           unitPrice,
			"quantity":             itemReq.Quantity,
			"line_total":           lineTotal,
			"discount":             itemDiscountAmount,
			"final_line_total":     finalLineTotal,
			"applied_discount_ids": appliedDiscountIDs,
		})

		// Create order item with server-calculated pricing
		orderItem := models.OrderItem{
			ProductID:  productID,
			VariantID:  orderVariantID,
			Name:       productName,
			Quantity:   itemReq.Quantity,
			UnitPrice:  unitPrice,
			TotalPrice: finalLineTotal, // Use server-calculated discounted price
			Image:      productImage,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Calculate final totals server-side
	finalTotal := subtotal - totalDiscount
	if finalTotal < 0 {
		finalTotal = 0
	}

	// Prepare applied discount IDs
	appliedDiscountIDs := make([]primitive.ObjectID, 0, len(appliedDiscountIDsMap))
	for id := range appliedDiscountIDsMap {
		appliedDiscountIDs = append(appliedDiscountIDs, id)
	}

	// Create order with server-calculated totals
	order := &models.Order{
		ShopID:             shopID,
		CustomerID:         customerID,
		Items:              orderItems,
		Subtotal:           subtotal,
		DiscountTotal:      totalDiscount,
		Total:              finalTotal,
		Status:             "pending",
		AppliedDiscountIDs: appliedDiscountIDs,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Save order to database
	created, err := CreateOrderService(order)
	if err != nil {
		return nil, err
	}

	// Record discount usage
	for id := range appliedDiscountIDsMap {
		for _, item := range itemDiscountDetails {
			for _, appliedID := range item["applied_discount_ids"].([]primitive.ObjectID) {
				if appliedID == id {
					_ = RecordDiscountUsageAtomic(id.Hex(), customerID, item["line_total"].(float64))
				}
			}
		}
	}

	// Reduce stock for all items
	for _, itemReq := range req.Items {
		productID, _ := primitive.ObjectIDFromHex(itemReq.ProductID)
		if itemReq.VariantID != "" {
			variantID, _ := primitive.ObjectIDFromHex(itemReq.VariantID)
			if !variantID.IsZero() {
				_ = ReduceVariantStock(productID, variantID, itemReq.Quantity)
			} else {
				_ = ReduceProductStock(productID, itemReq.Quantity)
			}
		} else {
			_ = ReduceProductStock(productID, itemReq.Quantity)
		}
	}

	return &OrderPlacementResult{
		Order:               created,
		ItemDiscountDetails: itemDiscountDetails,
	}, nil
}
