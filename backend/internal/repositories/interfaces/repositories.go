package interfaces

import (
	"context"
	"time"

	"shop24/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	FindByID(ctx context.Context, id string) (models.User, error)
}

type StoreRepository interface {
	Create(ctx context.Context, store models.Store) (models.Store, error)
	FindByOwner(ctx context.Context, ownerID string) (models.Store, error)
	FindBySlug(ctx context.Context, slug string) (models.Store, error)
	UpdateProfile(ctx context.Context, store models.Store) error
	SetPublished(ctx context.Context, ownerID, storeID string, published bool, at time.Time) error
	UpdateDraft(ctx context.Context, ownerID, storeID string, draft models.StorefrontConfig, note *models.DesignNote, at time.Time) error
	PublishDraft(ctx context.Context, ownerID, storeID string, at time.Time) error
}

type ProductFilter struct {
	Query      string
	Category   string
	Published  *bool
	LowStock   bool
	LowStockAt int
	Limit      int
}

type ProductRepository interface {
	Create(ctx context.Context, product models.Product) (models.Product, error)
	Update(ctx context.Context, product models.Product) (models.Product, error)
	Delete(ctx context.Context, storeID, id string) error
	FindByID(ctx context.Context, storeID, id string) (models.Product, error)
	List(ctx context.Context, storeID string, filter ProductFilter) ([]models.Product, error)
	AddStock(ctx context.Context, storeID, id string, delta int) error
	SetStock(ctx context.Context, storeID, id string, stock int, at time.Time) error
	Totals(ctx context.Context, storeID string) (products int, units int, err error)
}

type OrderFilter struct {
	Status        string
	ExcludeStatus string
	Since         time.Time
	Limit         int
}

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) (models.Order, error)
	FindByID(ctx context.Context, storeID, id string) (models.Order, error)
	List(ctx context.Context, storeID string, filter OrderFilter) ([]models.Order, error)
	UpdateStatus(ctx context.Context, storeID, id string, status models.OrderStatus, at time.Time) (models.Order, error)
	Summary(ctx context.Context, storeID string, since time.Time, exclude models.OrderStatus) (revenueCents int64, orders int, err error)
	RevenueByDay(ctx context.Context, storeID string, since time.Time, exclude models.OrderStatus) ([]models.RevenuePoint, error)
	BestSellers(ctx context.Context, storeID string, since time.Time, exclude models.OrderStatus, limit int) ([]models.BestSeller, error)
}
