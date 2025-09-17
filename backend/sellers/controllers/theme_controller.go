package controllers

import (
	"context"
	"net/http"

	"github.com/Endale2/DRPS/sellers/services"
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

	// require storeId from authenticated seller or payload
	storeIdStr := c.Param("storeId")
	if storeIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "storeId required"})
		return
	}

	sid, err := primitive.ObjectIDFromHex(storeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid storeId"})
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
	storeIdStr := c.Param("storeId")
	sid, err := primitive.ObjectIDFromHex(storeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid storeId"})
		return
	}
	themes, err := themeService.ListByStore(context.Background(), sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, themes)
}

func ActivateTheme(c *gin.Context) {
	storeIdStr := c.Param("storeId")
	themeIdStr := c.Param("themeId")
	sid, err := primitive.ObjectIDFromHex(storeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid storeId"})
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
