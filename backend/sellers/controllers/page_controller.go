package controllers

import (
	"context"
	"net/http"

	"github.com/Endale2/DRPS/sellers/services"
	sharedServices "github.com/Endale2/DRPS/shared/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var pageService = services.NewPageService()

// ensure packages are referenced (avoid unused import errors in some build contexts)
var _ = services.NewPageService
var _ = gin.Default

func CreatePage(c *gin.Context) {
	var p services.PageModelFromJSON
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Convert store id param
	storeIdStr := c.Param("storeId")
	sid, err := primitive.ObjectIDFromHex(storeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid storeId"})
		return
	}

	page := p.ToModel(sid)
	if err := pageService.Create(context.Background(), page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, page)
}

func UpdatePage(c *gin.Context) {
	pageIdStr := c.Param("pageId")
	pid, err := primitive.ObjectIDFromHex(pageIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageId"})
		return
	}
	var p services.PageModelFromJSON
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page := p.ToModel(primitive.NilObjectID)
	if err := pageService.Update(context.Background(), pid, page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func PublishPage(c *gin.Context) {
	pageIdStr := c.Param("pageId")
	pid, err := primitive.ObjectIDFromHex(pageIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageId"})
		return
	}
	// Fetch page first (reuse repository via service)
	// For simplicity, expect full page in body to validate and snapshot
	var p services.PageModelFromJSON
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page := p.ToModel(primitive.NilObjectID)
	if err := pageService.Publish(context.Background(), pid, page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ListPages returns all pages for a shop
func ListPages(c *gin.Context) {
	storeIdStr := c.Param("shopId")
	if storeIdStr == "" {
		storeIdStr = c.Param("storeId")
	}
	sid, err := primitive.ObjectIDFromHex(storeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId/storeId"})
		return
	}

	pages, err := pageService.ListByStore(context.Background(), sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pages": pages})
}

// GetPage returns a single page by ID
func GetPage(c *gin.Context) {
	pageIdStr := c.Param("pageId")
	pid, err := primitive.ObjectIDFromHex(pageIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageId"})
		return
	}

	page, err := pageService.GetByID(context.Background(), pid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}
	c.JSON(http.StatusOK, page)
}

// DeletePage removes a page by ID
func DeletePage(c *gin.Context) {
	pageIdStr := c.Param("pageId")
	pid, err := primitive.ObjectIDFromHex(pageIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageId"})
		return
	}

	if err := pageService.Delete(context.Background(), pid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Storefront endpoint: returns published page snapshot for a store and path
func GetStorefrontPage(c *gin.Context) {
	// Try to resolve storeId from query param first. If missing, resolve from the path param shopSlug.
	storeIdHex := c.Query("storeId")
	var sid primitive.ObjectID
	if storeIdHex != "" {
		// if provided, use it
		var err error
		sid, err = primitive.ObjectIDFromHex(storeIdHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid storeId"})
			return
		}
	} else {
		// fallback: resolve by slug (from path param)
		shopSlug := c.Param("shopSlug")
		if shopSlug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "storeId or shopSlug required"})
			return
		}
		shop, err := sharedServices.GetShopBySlugService(shopSlug)
		if err != nil || shop == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
			return
		}
		sid = shop.ID
	}
	path := c.Query("path")
	if path == "" {
		path = "/"
	}
	page, err := pageService.GetPublishedByStoreAndPath(context.Background(), sid, path)
	if err != nil {
		// If requesting Root and not found, try to create default home page
		if path == "/" {
			defaultPage, createErr := pageService.CreateDefaultHomePage(context.Background(), sid)
			if createErr == nil {
				c.JSON(http.StatusOK, gin.H{"snapshot": defaultPage.Snapshot})
				return
			}
			// If creation failed, log it but return original error (or createErr)
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"snapshot": page.Snapshot})
}
