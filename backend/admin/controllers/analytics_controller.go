package controllers

import (
	"net/http"

	"github.com/Endale2/DRPS/shared/repositories"
	"github.com/Endale2/DRPS/shared/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTotalSales returns the total number of sales (orders)
func GetTotalSales(c *gin.Context) {
	count, err := repositories.CountOrders(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting sales"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_sales": count})
}

// GetTotalRevenue returns the total revenue from all orders
func GetTotalRevenue(c *gin.Context) {
	orders, err := repositories.ListOrders(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
		return
	}

	totalRevenue := services.CalculateTotalRevenue(orders)
	c.JSON(http.StatusOK, gin.H{"total_revenue": totalRevenue})
}

// GetOrderCount returns the total number of orders
func GetOrderCount(c *gin.Context) {
	count, err := repositories.CountOrders(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_count": count})
}

// GetTopProducts returns the top N products by sales count
func GetTopProducts(c *gin.Context) {
	orders, err := repositories.ListOrders(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
		return
	}

	topProducts := services.GetTopProductsBySalesCount(orders, 5)
	c.JSON(http.StatusOK, topProducts)
}

// GetTopCustomers returns the top N customers by order count
func GetTopCustomers(c *gin.Context) {
	orders, err := repositories.ListOrders(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
		return
	}

	topCustomers := services.GetTopCustomersByOrderCountAsMap(orders, 5)
	c.JSON(http.StatusOK, topCustomers)
}
