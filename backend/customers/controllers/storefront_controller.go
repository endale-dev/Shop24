// customers/controllers/storefront_controller.go
package controllers

import (
	"net/http"

	sellerServices "github.com/Endale2/DRPS/sellers/services"
	shopServices "github.com/Endale2/DRPS/shared/services"
	"github.com/gin-gonic/gin"
)

// GetStorefront handles GET /shops/:shopSlug
func GetStorefront(c *gin.Context) {
	shopSlug := c.Param("shopSlug")

	// 1. Fetch Shop
	shop, err := shopServices.GetShopBySlugService(shopSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// 2. Fetch Active Theme
	themeService := sellerServices.NewThemeService()
	theme, err := themeService.GetActiveThemeByShopSlug(c.Request.Context(), shopSlug)
	if err != nil || theme == nil {
		// No active theme found in DB. Create and persist a default one.
		t, err := themeService.CreateDefaultTheme(c.Request.Context(), shop.ID)
		if err != nil {
			// Fallback to transient default only if creation fails (DB error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize theme"})
			return
		}
		// Convert pointer to sharedModels.Theme if CreateDefaultTheme returns something else,
		// but check theme_service signature. It returns *sharedModels.Theme.
		theme = t
	}

	// 3. Construct Response
	response := gin.H{
		"shop":  shop,
		"theme": theme,
		// Add default configs for other sections expected by frontend
		"layout": gin.H{
			"headerStyle":    theme.Layout.HeaderStyle,
			"footerStyle":    theme.Layout.HeaderStyle, // fallback
			"containerWidth": theme.Layout.ContainerWidth,
			"gridColumns":    theme.Layout.ProductGridCols,
			"sidebarEnabled": false,
		},
		"headerBlocks": theme.HeaderBlocks, // Explicit layout
		"footerBlocks": theme.FooterBlocks, // Explicit layout
		"seo": gin.H{
			"title":       shop.Name,
			"description": shop.Description,
		},
		"navigation": gin.H{
			"style": "scrolled",
		},
		"footer": gin.H{
			"style": "simple",
		},
	}

	c.JSON(http.StatusOK, response)
}
