package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"shop24/internal/apperr"
	"shop24/internal/domain/models"
	"shop24/internal/httpx"
	"shop24/internal/repositories/interfaces"
	"shop24/internal/services"
)

type StoreHandler struct {
	stores *services.StoreService
	log    *slog.Logger
}

func NewStoreHandler(stores *services.StoreService, log *slog.Logger) *StoreHandler {
	return &StoreHandler{stores: stores, log: log}
}

func (h *StoreHandler) Create(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	store, err := h.stores.Create(r.Context(), owner, services.StoreInput{Name: body.Name, Slug: body.Slug, Description: body.Description})
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusCreated, store)
}

func (h *StoreHandler) Mine(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	store, err := h.stores.Mine(r.Context(), owner)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, store)
}

func (h *StoreHandler) Update(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		LogoURL     string `json:"logoUrl"`
		Currency    string `json:"currency"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	store, err := h.stores.Update(r.Context(), owner, services.StoreInput{
		Name: body.Name, Slug: body.Slug, Description: body.Description, LogoURL: body.LogoURL, Currency: body.Currency,
	})
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, store)
}

func (h *StoreHandler) Publish(w http.ResponseWriter, r *http.Request) {
	h.setPublished(w, r, true)
}

func (h *StoreHandler) Unpublish(w http.ResponseWriter, r *http.Request) {
	h.setPublished(w, r, false)
}

func (h *StoreHandler) setPublished(w http.ResponseWriter, r *http.Request, published bool) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	store, err := h.stores.SetPublished(r.Context(), owner, published)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, store)
}

type ProductHandler struct {
	products *services.ProductService
	log      *slog.Logger
}

func NewProductHandler(products *services.ProductService, log *slog.Logger) *ProductHandler {
	return &ProductHandler{products: products, log: log}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	body, err := productBody(w, r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	product, err := h.products.Create(r.Context(), owner, body)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusCreated, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	body, err := productBody(w, r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	product, err := h.products.Update(r.Context(), owner, r.PathValue("id"), body)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, product)
}

func (h *ProductHandler) SetStock(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Stock int `json:"stock"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	product, err := h.products.SetStock(r.Context(), owner, r.PathValue("id"), body.Stock)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	if err := h.products.Delete(r.Context(), owner, r.PathValue("id")); err != nil {
		fail(h.log, w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	product, err := h.products.Get(r.Context(), owner, r.PathValue("id"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	filter := interfaces.ProductFilter{
		Query: r.URL.Query().Get("q"), Category: r.URL.Query().Get("category"), Limit: 100,
	}
	switch r.URL.Query().Get("status") {
	case "", "all":
	case "published":
		published := true
		filter.Published = &published
	case "draft":
		published := false
		filter.Published = &published
	default:
		fail(h.log, w, apperr.New(http.StatusBadRequest, "invalid_filter", "Status must be all, published, or draft"))
		return
	}
	if r.URL.Query().Get("lowStock") == "1" {
		filter.LowStock = true
		filter.LowStockAt = services.LowStockThreshold
	}
	products, err := h.products.List(r.Context(), owner, filter)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, map[string]any{"items": products})
}

func productBody(w http.ResponseWriter, r *http.Request) (services.ProductInput, error) {
	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		PriceCents  int64    `json:"priceCents"`
		Images      []string `json:"images"`
		Category    string   `json:"category"`
		SKU         string   `json:"sku"`
		Stock       int      `json:"stock"`
		Published   bool     `json:"published"`
	}
	if err := decode(w, r, &body); err != nil {
		return services.ProductInput{}, err
	}
	return services.ProductInput{
		Name: body.Name, Description: body.Description, PriceCents: body.PriceCents, Images: body.Images,
		Category: body.Category, SKU: body.SKU, Stock: body.Stock, Published: body.Published,
	}, nil
}

type OrderHandler struct {
	orders *services.OrderService
	log    *slog.Logger
}

func NewOrderHandler(orders *services.OrderService, log *slog.Logger) *OrderHandler {
	return &OrderHandler{orders: orders, log: log}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	orders, err := h.orders.List(r.Context(), owner, r.URL.Query().Get("status"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, map[string]any{"items": orders})
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	order, err := h.orders.Get(r.Context(), owner, r.PathValue("id"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	order, err := h.orders.UpdateStatus(r.Context(), owner, r.PathValue("id"), models.OrderStatus(body.Status))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, order)
}

type AnalyticsHandler struct {
	analytics *services.AnalyticsService
	log       *slog.Logger
}

func NewAnalyticsHandler(analytics *services.AnalyticsService, log *slog.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{analytics: analytics, log: log}
}

func (h *AnalyticsHandler) Overview(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	overview, err := h.analytics.Overview(r.Context(), owner, days)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, overview)
}
