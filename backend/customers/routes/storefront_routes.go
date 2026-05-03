package routes

import (
	"github.com/Endale2/DRPS/auth/middlewares"
	custControllers "github.com/Endale2/DRPS/customers/controllers"
	sellerControllers "github.com/Endale2/DRPS/sellers/controllers"
	"github.com/gin-gonic/gin"
)

func StorefrontRoutes(r *gin.Engine) {
	// All storefront endpoints begin with /shops/:shopSlug
	shops := r.Group("/shops")
	{
		// GET /shops/:shopSlug → view shop details (legacy)
		shops.GET("/:shopSlug", custControllers.GetStorefront)

		// GET /shops/:shopSlug/theme → get active theme (public)
		shops.GET("/:shopSlug/theme", custControllers.GetShopTheme)

		// Collection endpoints nested under /shops/:shopSlug/collections
		collections := shops.Group("/:shopSlug/collections")
		{
			// list all collections for this shop
			collections.GET("", custControllers.ListCollections)
			// get one collection + its product summaries
			collections.GET("/:collectionHandle", custControllers.GetCollection)
		}

		// Product endpoints nested under /shops/:shopSlug/products
		products := shops.Group("/:shopSlug/products")
		{
			products.GET("", custControllers.GetProductsByShop)
			products.GET("search", custControllers.SearchProductsByShop)
			products.GET("/:productSlug", custControllers.GetProductDetail)
			products.GET("id/:productId", custControllers.GetProductDetailByID)
		}

		// Debug endpoints (no auth required for testing)
		shops.GET("/:shopSlug/debug/product/:productId", custControllers.DebugProduct)
		// Page rendering endpoint (public) - returns published snapshot for a path
		shops.GET("/:shopSlug/page", func(c *gin.Context) {
			// delegate to sellers controllers page endpoint to avoid dependency cycles
			// expect query param storeId and path
			sellerControllers.GetStorefrontPage(c)
		})
		shops.POST("/:shopSlug/debug/product/:productId/price", custControllers.UpdateProductPrice)
		shops.POST("/:shopSlug/debug/fix-variant-ids", custControllers.FixVariantIDs)
		shops.GET("/:shopSlug/test-discounts", custControllers.TestDiscounts)

		// Protected endpoints: all under /shops/:shopSlug/ and require auth
		auth := shops.Group("/:shopSlug", middlewares.AuthMiddleware())
		{
			auth.GET("/cart", custControllers.GetCart)
			auth.POST("/cart/items", custControllers.AddToCart)
			auth.PUT("/cart/items", custControllers.UpdateCartItem)
			auth.DELETE("/cart/items", custControllers.RemoveCartItem)
			auth.POST("/cart/clear", custControllers.ClearCart)
			auth.POST("/orders", custControllers.PlaceOrder)
			auth.GET("/orders", custControllers.ListShopOrders)
			auth.GET("/orders/:orderId", custControllers.GetOrderDetail)
			auth.GET("/wishlist", custControllers.GetWishlist)
			auth.POST("/wishlist", custControllers.AddToWishlist)
			auth.DELETE("/wishlist/:productId", custControllers.RemoveFromWishlist)
		}
	}
}
