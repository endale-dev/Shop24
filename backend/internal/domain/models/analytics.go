package models

type RevenuePoint struct {
	Date         string `json:"date"`
	RevenueCents int64  `json:"revenueCents"`
	Orders       int    `json:"orders"`
}

type BestSeller struct {
	ProductID    string `json:"productId"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	RevenueCents int64  `json:"revenueCents"`
}

type Overview struct {
	RangeDays              int            `json:"rangeDays"`
	TotalRevenueCents      int64          `json:"totalRevenueCents"`
	TotalOrders            int            `json:"totalOrders"`
	AverageOrderValueCents int64          `json:"averageOrderValueCents"`
	TotalProducts          int            `json:"totalProducts"`
	UnitsInStock           int            `json:"unitsInStock"`
	LowStockThreshold      int            `json:"lowStockThreshold"`
	LowStock               []Product      `json:"lowStock"`
	RecentOrders           []Order        `json:"recentOrders"`
	BestSellers            []BestSeller   `json:"bestSellers"`
	RevenueOverTime        []RevenuePoint `json:"revenueOverTime"`
}
