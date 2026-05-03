// customers/controllers/order_controller.go
package controllers

import (
	"net/http"

	"github.com/Endale2/DRPS/shared/models"
	"github.com/Endale2/DRPS/shared/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DebugProduct handles GET /shops/:shopSlug/debug/product/:productId
func DebugProduct(c *gin.Context) {
	shopSlug := c.Param("shopSlug")
	productID := c.Param("productId")

	// Lookup the shop by its slug
	shop, err := services.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup shop"})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Parse product ID
	pid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Get customer ID from context if available
	var customerID *primitive.ObjectID
	if cidVal, exists := c.Get("user_id"); exists {
		if cidHex, ok := cidVal.(string); ok {
			if cid, err := primitive.ObjectIDFromHex(cidHex); err == nil {
				customerID = &cid
			}
		}
	}

	// Fetch product with discounts applied
	product, err := services.GetProductByIDWithDiscountsService(pid.Hex(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product: " + err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Verify product belongs to this shop
	if product.ShopID != shop.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product does not belong to this shop"})
		return
	}

	// Build debug response
	variants := make([]map[string]interface{}, 0, len(product.Variants))
	for i, v := range product.Variants {
		variant := map[string]interface{}{
			"index":      i,
			"variant_id": v.VariantID.Hex(),
			"options":    v.Options,
			"price":      v.Price,
			"stock":      v.Stock,
			"image":      v.Image,
		}
		variants = append(variants, variant)
	}

	response := map[string]interface{}{
		"product_id":     product.ID.Hex(),
		"name":           product.Name,
		"shop_id":        product.ShopID.Hex(),
		"product_price":  product.Price,
		"product_stock":  product.Stock,
		"total_variants": len(product.Variants),
		"variants":       variants,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProductPrice handles POST /shops/:shopSlug/debug/product/:productId/price
// This is a temporary endpoint for testing - update product price
func UpdateProductPrice(c *gin.Context) {
	shopSlug := c.Param("shopSlug")
	productID := c.Param("productId")

	// Lookup the shop by its slug
	shop, err := services.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup shop"})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Parse product ID
	pid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Parse request body
	var req struct {
		Price float64 `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update product price
	_, err = services.UpdateProductService(pid.Hex(), bson.M{"price": req.Price})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product price updated successfully", "new_price": req.Price})
}

// FixVariantIDs handles POST /shops/:shopSlug/debug/fix-variant-ids
// This is a temporary endpoint for fixing variant IDs in existing products
func FixVariantIDs(c *gin.Context) {
	shopSlug := c.Param("shopSlug")

	// Lookup the shop by its slug
	shop, err := services.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup shop"})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Run the variant ID fix
	err = services.UpdateProductVariantIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fix variant IDs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variant IDs fixed successfully"})
}

// TestDiscounts handles GET /shops/:shopSlug/test-discounts
func TestDiscounts(c *gin.Context) {
	shopSlug := c.Param("shopSlug")
	shop, err := services.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup shop"})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Get all discounts for this shop
	discounts, err := services.ListDiscountsByShopService(shop.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get customer ID if authenticated
	var customerID primitive.ObjectID
	var customerSegmentIDs []primitive.ObjectID
	if cidVal, exists := c.Get("user_id"); exists {
		if cidHex, ok := cidVal.(string); ok {
			if cid, err := primitive.ObjectIDFromHex(cidHex); err == nil {
				customerID = cid
				customerSegmentIDs, _ = services.GetCustomerSegmentIDs(shop.ID, customerID)
			}
		}
	}

	// Test each discount
	var results []map[string]interface{}
	for _, discount := range discounts {
		result := map[string]interface{}{
			"id":                  discount.ID.Hex(),
			"name":                discount.Name,
			"type":                discount.Type,
			"value":               discount.Value,
			"category":            discount.Category,
			"active":              services.IsDiscountActive(&discount),
			"shop_id":             discount.ShopID.Hex(),
			"applies_to_products": discount.AppliesToProducts,
			"applies_to_variants": discount.AppliesToVariants,
			"eligibility_type":    discount.EligibilityType,
			"allowed_customers":   discount.AllowedCustomerIDs,
			"allowed_segments":    discount.AllowedSegmentIDs,
			"usage_limit":         discount.UsageLimit,
			"per_customer_limit":  discount.PerCustomerLimit,
			"current_usage":       discount.CurrentUsage,
			"start_at":            discount.StartAt,
			"end_at":              discount.EndAt,
		}

		// Test customer eligibility
		if !customerID.IsZero() {
			canUse, err := services.CanCustomerUseDiscount(&discount, customerID, customerSegmentIDs)
			result["customer_can_use"] = canUse
			result["customer_error"] = err
		}

		// Test discount calculation
		testPrice := 100.0
		savings := services.CalculateDiscountAmount(&discount, testPrice)
		result["test_calculation"] = map[string]interface{}{
			"test_price":  testPrice,
			"savings":     savings,
			"final_price": testPrice - savings,
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"shop_id":           shop.ID.Hex(),
		"customer_id":       customerID.Hex(),
		"customer_segments": customerSegmentIDs,
		"total_discounts":   len(discounts),
		"discounts":         results,
	})
}

// PlaceOrder handles POST /shops/:shopSlug/orders
func PlaceOrder(c *gin.Context) {
	// Get shop by slug
	shopSlug := c.Param("shopSlug")
	shop, err := services.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup shop"})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Get authenticated customer ID
	cidVal, exists := c.Get("user_id")
	if !exists || cidVal == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	cidHex, ok := cidVal.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}
	customerID, err := primitive.ObjectIDFromHex(cidHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid customer ID"})
		return
	}

	// Link customer to shop if not already linked
	_, _, _ = services.LinkIfNotLinked(shop.ID, customerID)

	// Parse request
	var req services.PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no items in order"})
		return
	}

	// Delegate to service - ALL business logic handled there
	result, err := services.PlaceOrderFromRequest(shop.ID, customerID, req)
	if err != nil {
		// Map service errors to HTTP status codes
		switch err {
		case services.ErrInvalidProductID, services.ErrInvalidVariantID, services.ErrInvalidQuantity:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case services.ErrProductNotFound, services.ErrVariantNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case services.ErrProductNotInShop:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case services.ErrInsufficientStock:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"id":                    result.Order.ID.Hex(),
		"order":                 result.Order,
		"item_discount_details": result.ItemDiscountDetails,
		"server_calculated":     true,
		"security_note":         "All pricing calculated securely on server",
	})
}

// ListShopOrders handles GET /shops/:shopSlug/orders
func ListShopOrders(c *gin.Context) {
	shopSlug := c.Param("shopSlug")
	shop, err := services.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup shop"})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// List by shop ID
	list, err := services.ListOrdersByShopService(shop.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter to only the current customer, if logged in
	if cidVal, exists := c.Get("user_id"); exists {
		if cidHex, ok := cidVal.(string); ok {
			if cid, err := primitive.ObjectIDFromHex(cidHex); err == nil {
				filtered := make([]models.Order, 0, len(list))
				for _, o := range list {
					if o.CustomerID == cid {
						filtered = append(filtered, o)
					}
				}
				list = filtered
			}
		}
	}

	c.JSON(http.StatusOK, list)
}

// GetOrderDetail handles GET /shops/:shopSlug/orders/:orderId
func GetOrderDetail(c *gin.Context) {
	shopSlug := c.Param("shopSlug")
	shop, err := services.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup shop"})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	orderID := c.Param("orderId")
	order, err := services.GetOrderByIDService(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if order == nil || order.ShopID != shop.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Only the placing customer may view it
	if cidVal, exists := c.Get("user_id"); exists {
		if cidHex, ok := cidVal.(string); ok {
			if cid, err := primitive.ObjectIDFromHex(cidHex); err == nil {
				if order.CustomerID != cid {
					c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to view this order"})
					return
				}
			}
		}
	}

	c.JSON(http.StatusOK, order)
}
