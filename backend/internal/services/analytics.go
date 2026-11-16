package services

import (
	"context"
	"sort"
	"time"

	"shop24/internal/domain/models"
	"shop24/internal/repositories/interfaces"
)

type AnalyticsService struct {
	stores   interfaces.StoreRepository
	products interfaces.ProductRepository
	orders   interfaces.OrderRepository
}

func NewAnalyticsService(stores interfaces.StoreRepository, products interfaces.ProductRepository, orders interfaces.OrderRepository) *AnalyticsService {
	return &AnalyticsService{stores: stores, products: products, orders: orders}
}

func (s *AnalyticsService) Overview(ctx context.Context, ownerID string, days int) (models.Overview, error) {
	if days != 7 && days != 30 && days != 90 {
		days = 30
	}
	store, err := requireStore(ctx, s.stores, ownerID)
	if err != nil {
		return models.Overview{}, err
	}
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(days - 1))
	revenue, orderCount, err := s.orders.Summary(ctx, store.ID, start, models.OrderCancelled)
	if err != nil {
		return models.Overview{}, err
	}
	points, err := s.orders.RevenueByDay(ctx, store.ID, start, models.OrderCancelled)
	if err != nil {
		return models.Overview{}, err
	}
	best, err := s.orders.BestSellers(ctx, store.ID, start, models.OrderCancelled, 5)
	if err != nil {
		return models.Overview{}, err
	}
	recent, err := s.orders.List(ctx, store.ID, interfaces.OrderFilter{ExcludeStatus: string(models.OrderCancelled), Since: start, Limit: 5})
	if err != nil {
		return models.Overview{}, err
	}
	productCount, units, err := s.products.Totals(ctx, store.ID)
	if err != nil {
		return models.Overview{}, err
	}
	low, err := s.products.List(ctx, store.ID, interfaces.ProductFilter{LowStock: true, LowStockAt: LowStockThreshold, Limit: 8})
	if err != nil {
		return models.Overview{}, err
	}
	sort.Slice(low, func(i, j int) bool { return low[i].Stock < low[j].Stock })
	average := int64(0)
	if orderCount > 0 {
		average = revenue / int64(orderCount)
	}
	if best == nil {
		best = []models.BestSeller{}
	}
	if recent == nil {
		recent = []models.Order{}
	}
	return models.Overview{
		RangeDays: days, TotalRevenueCents: revenue, TotalOrders: orderCount, AverageOrderValueCents: average,
		TotalProducts: productCount, UnitsInStock: units, LowStockThreshold: LowStockThreshold,
		LowStock: presentProducts(low), RecentOrders: recent, BestSellers: best,
		RevenueOverTime: fillDays(start, days, points),
	}, nil
}

func fillDays(start time.Time, days int, points []models.RevenuePoint) []models.RevenuePoint {
	byDate := map[string]models.RevenuePoint{}
	for _, point := range points {
		byDate[point.Date] = point
	}
	out := make([]models.RevenuePoint, 0, days)
	for i := 0; i < days; i++ {
		date := start.AddDate(0, 0, i).Format("2006-01-02")
		point := byDate[date]
		point.Date = date
		out = append(out, point)
	}
	return out
}
