package services

import (
	"sort"
	"time"

	"github.com/Endale2/DRPS/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============================================================================
// ANALYTICS BUSINESS LOGIC
// Extracted from admin/controllers/analytics_controller.go
// and sellers/controllers/seller_analytics_controller.go
// ============================================================================

// ProductSalesStats represents sales statistics for a product
type ProductSalesStats struct {
	ProductID    primitive.ObjectID `json:"product_id"`
	ProductName  string             `json:"product_name"`
	Category     string             `json:"category,omitempty"`
	MainImage    string             `json:"main_image,omitempty"`
	TotalSold    int                `json:"total_sold"`
	TotalRevenue float64            `json:"total_revenue"`
}

// CustomerOrderStats represents order statistics for a customer
type CustomerOrderStats struct {
	CustomerID primitive.ObjectID `json:"customer_id"`
	OrderCount int                `json:"order_count"`
}

// TimeSeriesDataPoint represents a single point in time-series data
type TimeSeriesDataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// ============================================================================
// ORDER FILTERING
// ============================================================================

// FilterSoldOrders filters orders to only include those with "sold" status
// (paid, shipped, or delivered)
func FilterSoldOrders(orders []models.Order) []models.Order {
	var filtered []models.Order
	for _, o := range orders {
		status := o.Status
		if status == "paid" || status == "shipped" || status == "delivered" {
			filtered = append(filtered, o)
		}
	}
	return filtered
}

// ============================================================================
// REVENUE CALCULATIONS
// ============================================================================

// CalculateTotalRevenue calculates the total revenue from a list of orders
func CalculateTotalRevenue(orders []models.Order) float64 {
	total := 0.0
	for _, o := range orders {
		total += o.Total
	}
	return total
}

// CalculateAverageOrderValue calculates the average order value
func CalculateAverageOrderValue(orders []models.Order) float64 {
	if len(orders) == 0 {
		return 0.0
	}
	totalRevenue := CalculateTotalRevenue(orders)
	return totalRevenue / float64(len(orders))
}

// ============================================================================
// PRODUCT ANALYTICS
// ============================================================================

// GetTopProductsByRevenue returns the top N products ranked by total revenue
func GetTopProductsByRevenue(orders []models.Order, products []models.Product, limit int) []ProductSalesStats {
	// Calculate sales per product
	productSales := make(map[primitive.ObjectID]int)
	productRevenue := make(map[primitive.ObjectID]float64)

	for _, o := range orders {
		for _, item := range o.Items {
			productSales[item.ProductID] += item.Quantity
			productRevenue[item.ProductID] += item.TotalPrice
		}
	}

	// Build stats with product details
	var stats []ProductSalesStats
	for _, p := range products {
		stats = append(stats, ProductSalesStats{
			ProductID:    p.ID,
			ProductName:  p.Name,
			MainImage:    p.MainImage,
			TotalSold:    productSales[p.ID],
			TotalRevenue: productRevenue[p.ID],
		})
	}

	// Sort by revenue descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalRevenue > stats[j].TotalRevenue
	})

	// Limit results
	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats
}

// GetTopProductsBySalesCount returns top products by quantity sold (for admin analytics)
func GetTopProductsBySalesCount(orders []models.Order, limit int) []struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
} {
	productSales := make(map[string]int)
	for _, o := range orders {
		for _, item := range o.Items {
			productSales[item.ProductID.Hex()] += item.Quantity
		}
	}

	// Convert to slice for sorting
	type kv struct {
		ProductID string
		Quantity  int
	}
	var sorted []kv
	for k, v := range productSales {
		sorted = append(sorted, kv{k, v})
	}

	// Sort by quantity descending
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Quantity > sorted[j].Quantity
	})

	// Limit results
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}

	// Convert to return type
	result := make([]struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}, len(sorted))
	for i, kv := range sorted {
		result[i].ProductID = kv.ProductID
		result[i].Quantity = kv.Quantity
	}

	return result
}

// ============================================================================
// CUSTOMER ANALYTICS
// ============================================================================

// GetTopCustomersByOrderCount returns the top N customers by number of orders
func GetTopCustomersByOrderCount(orders []models.Order, limit int) []CustomerOrderStats {
	customerOrders := make(map[primitive.ObjectID]int)
	for _, o := range orders {
		customerOrders[o.CustomerID]++
	}

	// Convert to slice for sorting
	var stats []CustomerOrderStats
	for customerID, count := range customerOrders {
		stats = append(stats, CustomerOrderStats{
			CustomerID: customerID,
			OrderCount: count,
		})
	}

	// Sort by order count descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].OrderCount > stats[j].OrderCount
	})

	// Limit results
	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats
}

// GetTopCustomersByOrderCountAsMap returns top customers for admin (using hex IDs)
func GetTopCustomersByOrderCountAsMap(orders []models.Order, limit int) []struct {
	CustomerID string `json:"customer_id"`
	OrderCount int    `json:"order_count"`
} {
	customerOrders := make(map[string]int)
	for _, o := range orders {
		customerOrders[o.CustomerID.Hex()]++
	}

	type kv struct {
		CustomerID string
		OrderCount int
	}
	var sorted []kv
	for k, v := range customerOrders {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OrderCount > sorted[j].OrderCount
	})

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}

	result := make([]struct {
		CustomerID string `json:"customer_id"`
		OrderCount int    `json:"order_count"`
	}, len(sorted))
	for i, kv := range sorted {
		result[i].CustomerID = kv.CustomerID
		result[i].OrderCount = kv.OrderCount
	}

	return result
}

// CountNewCustomers counts customers created within the last N days
func CountNewCustomers(customers []models.ShopCustomer, days int) int {
	cutoff := time.Now().AddDate(0, 0, -days)
	count := 0
	for _, cst := range customers {
		if cst.CreatedAt.After(cutoff) {
			count++
		}
	}
	return count
}

// CountReturningCustomers counts customers who have placed more than one order
func CountReturningCustomers(customers []models.ShopCustomer, orders []models.Order) int {
	customerOrderCount := make(map[primitive.ObjectID]int)
	for _, o := range orders {
		customerOrderCount[o.CustomerID]++
	}

	returningCount := 0
	for _, cst := range customers {
		if customerOrderCount[cst.CustomerID] > 1 {
			returningCount++
		}
	}
	return returningCount
}

// ============================================================================
// TIME SERIES ANALYTICS
// ============================================================================

// GetRevenueOverTime calculates daily revenue for the last N days
func GetRevenueOverTime(orders []models.Order, days int) (labels []string, values []float64) {
	revenueMap := make(map[string]float64)
	now := time.Now()

	// Initialize all days with 0
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		revenueMap[date] = 0
	}

	// Aggregate revenue by date
	for _, o := range orders {
		date := o.CreatedAt.Format("2006-01-02")
		if _, exists := revenueMap[date]; exists {
			revenueMap[date] += o.Total
		}
	}

	// Build result arrays in chronological order
	labels = make([]string, days)
	values = make([]float64, days)
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		labels[days-1-i] = date
		values[days-1-i] = revenueMap[date]
	}

	return labels, values
}

// GetOrdersOverTime calculates daily order count for the last N days
func GetOrdersOverTime(orders []models.Order, days int) (labels []string, values []int) {
	orderMap := make(map[string]int)
	now := time.Now()

	// Initialize all days with 0
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		orderMap[date] = 0
	}

	// Count orders by date
	for _, o := range orders {
		date := o.CreatedAt.Format("2006-01-02")
		if _, exists := orderMap[date]; exists {
			orderMap[date]++
		}
	}

	// Build result arrays in chronological order
	labels = make([]string, days)
	values = make([]int, days)
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		labels[days-1-i] = date
		values[days-1-i] = orderMap[date]
	}

	return labels, values
}

// GetCustomersOverTime calculates daily customer registration for the last N days
func GetCustomersOverTime(customers []models.ShopCustomer, days int) (labels []string, values []int) {
	customerMap := make(map[string]int)
	now := time.Now()

	// Initialize all days with 0
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		customerMap[date] = 0
	}

	// Count customers by registration date
	for _, cst := range customers {
		date := cst.CreatedAt.Format("2006-01-02")
		if _, exists := customerMap[date]; exists {
			customerMap[date]++
		}
	}

	// Build result arrays in chronological order
	labels = make([]string, days)
	values = make([]int, days)
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		labels[days-1-i] = date
		values[days-1-i] = customerMap[date]
	}

	return labels, values
}

// ============================================================================
// CATEGORY/COLLECTION ANALYTICS
// ============================================================================

// CategorySales represents sales data for a category
type CategorySales struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
}

// GetCategorySales calculates revenue by category for the last N days
func GetCategorySales(orders []models.Order, productCategories map[primitive.ObjectID]string, days int) []CategorySales {
	categoryRevenue := make(map[string]float64)
	now := time.Now()
	cutoff := now.AddDate(0, 0, -days)

	// Aggregate revenue by category
	for _, o := range orders {
		if o.CreatedAt.Before(cutoff) {
			continue
		}
		for _, item := range o.Items {
			category := productCategories[item.ProductID]
			if category == "" {
				category = "Uncategorized"
			}
			categoryRevenue[category] += item.TotalPrice
		}
	}

	// Convert to slice and sort
	var pairs []CategorySales
	for cat, rev := range categoryRevenue {
		pairs = append(pairs, CategorySales{
			Category: cat,
			Revenue:  rev,
		})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Revenue > pairs[j].Revenue
	})

	// Limit to top 6
	if len(pairs) > 6 {
		pairs = pairs[:6]
	}

	return pairs
}

// ============================================================================
// INVENTORY ANALYTICS
// ============================================================================

// InventoryStatus represents the inventory status summary
type InventoryStatus struct {
	TotalProducts int `json:"total_products"`
	InStock       int `json:"in_stock"`
	LowStock      int `json:"low_stock"`
	OutOfStock    int `json:"out_of_stock"`
}

// GetInventoryStatus calculates inventory statistics
func GetInventoryStatus(products []models.Product) InventoryStatus {
	status := InventoryStatus{
		TotalProducts: len(products),
	}

	for _, p := range products {
		if p.Stock > 10 {
			status.InStock++
		} else if p.Stock > 0 {
			status.LowStock++
		} else {
			status.OutOfStock++
		}
	}

	return status
}

// CountLowStockProducts counts products with stock below threshold
func CountLowStockProducts(products []models.Product, threshold int) int {
	count := 0
	for _, p := range products {
		if p.Stock < threshold {
			count++
		}
	}
	return count
}
