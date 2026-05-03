package controllers

import (
	"context"
	"net/http"

	"github.com/Endale2/DRPS/sellers/services"
	sharedModels "github.com/Endale2/DRPS/shared/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var themeService = services.NewThemeService()

func CreateTheme(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// require shopId from authenticated seller or payload
	shopIdStr := c.Param("shopId")
	if shopIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shopId required"})
		return
	}

	sid, err := primitive.ObjectIDFromHex(shopIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId"})
		return
	}

	t := services.ThemeModelFromMap(sid, payload)
	if err := themeService.Create(context.Background(), t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func ListThemes(c *gin.Context) {
	shopIdStr := c.Param("shopId")
	sid, err := primitive.ObjectIDFromHex(shopIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId"})
		return
	}
	themes, err := themeService.ListByStore(context.Background(), sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, themes)
}

func GetActiveTheme(c *gin.Context) {
	shopIdStr := c.Param("shopId")
	sid, err := primitive.ObjectIDFromHex(shopIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId"})
		return
	}

	theme, err := themeService.GetActiveThemeByShopID(context.Background(), sid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active theme found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"theme": theme})
}

func ActivateTheme(c *gin.Context) {
	shopIdStr := c.Param("shopId")
	themeIdStr := c.Param("themeId")
	sid, err := primitive.ObjectIDFromHex(shopIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId"})
		return
	}
	tid, err := primitive.ObjectIDFromHex(themeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid themeId"})
		return
	}
	if err := themeService.Activate(context.Background(), sid, tid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func UpdateThemeColors(c *gin.Context) {
	themeIdStr := c.Param("themeId")
	tid, err := primitive.ObjectIDFromHex(themeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid themeId"})
		return
	}

	var colors sharedModels.ThemeColors
	if err := c.BindJSON(&colors); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := themeService.UpdateThemeColors(context.Background(), tid, colors); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "colors": colors})
}

func UpdateThemeTypography(c *gin.Context) {
	themeIdStr := c.Param("themeId")
	tid, err := primitive.ObjectIDFromHex(themeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid themeId"})
		return
	}

	var typography sharedModels.ThemeTypography
	if err := c.BindJSON(&typography); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := themeService.UpdateThemeTypography(context.Background(), tid, typography); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "typography": typography})
}

func UpdateThemeLayout(c *gin.Context) {
	themeIdStr := c.Param("themeId")
	tid, err := primitive.ObjectIDFromHex(themeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid themeId"})
		return
	}

	var layout sharedModels.ThemeLayout
	if err := c.BindJSON(&layout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := themeService.UpdateThemeLayout(context.Background(), tid, layout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "layout": layout})
}

func UpdateThemeAssets(c *gin.Context) {
	themeIdStr := c.Param("themeId")
	tid, err := primitive.ObjectIDFromHex(themeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid themeId"})
		return
	}

	var assets sharedModels.ThemeAssets
	if err := c.BindJSON(&assets); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := themeService.UpdateThemeAssets(context.Background(), tid, assets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "assets": assets})
}
