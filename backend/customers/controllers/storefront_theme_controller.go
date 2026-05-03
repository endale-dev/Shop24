// customers/controllers/storefront_theme_controller.go
package controllers

import (
	"context"
	"net/http"

	"github.com/Endale2/DRPS/sellers/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var themeService = services.NewThemeService()

// GetShopTheme returns the active theme for a shop (public endpoint for storefront)
// GET /shops/:shopSlug/theme
func GetShopTheme(c *gin.Context) {
	shopSlug := c.Param("shopSlug")
	if shopSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shop slug is required"})
		return
	}

	ctx := context.Background()
	theme, err := themeService.GetActiveThemeByShopSlug(ctx, shopSlug)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No theme found, return default
			c.JSON(http.StatusOK, gin.H{
				"theme":   nil,
				"message": "No active theme found, using defaults",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch theme"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"theme": theme})
}
