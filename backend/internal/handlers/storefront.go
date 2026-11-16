package handlers

import (
	"log/slog"
	"net/http"

	"shop24/internal/domain/models"
	"shop24/internal/httpx"
	"shop24/internal/services"
)

type StorefrontHandler struct {
	designs *services.StorefrontService
	builder *services.AIBuilderService
	log     *slog.Logger
}

func NewStorefrontHandler(designs *services.StorefrontService, builder *services.AIBuilderService, log *slog.Logger) *StorefrontHandler {
	return &StorefrontHandler{designs: designs, builder: builder, log: log}
}

func (h *StorefrontHandler) Generate(w http.ResponseWriter, r *http.Request) {
	h.prompt(w, r, false)
}

func (h *StorefrontHandler) Refine(w http.ResponseWriter, r *http.Request) {
	h.prompt(w, r, true)
}

func (h *StorefrontHandler) prompt(w http.ResponseWriter, r *http.Request, refine bool) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Prompt string `json:"prompt"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	var result services.DesignResult
	if refine {
		result, err = h.builder.Refine(r.Context(), owner, body.Prompt)
	} else {
		result, err = h.builder.Generate(r.Context(), owner, body.Prompt)
	}
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, result)
}

func (h *StorefrontHandler) Save(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	var body struct {
		Config models.StorefrontConfig `json:"config"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	store, err := h.designs.Save(r.Context(), owner, body.Config)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, store)
}

func (h *StorefrontHandler) Publish(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	store, err := h.designs.Publish(r.Context(), owner)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, store)
}

type PublicHandler struct {
	stores      *services.StoreService
	products    *services.ProductService
	orders      *services.OrderService
	collections *services.CollectionService
	log         *slog.Logger
}

func NewPublicHandler(stores *services.StoreService, products *services.ProductService, orders *services.OrderService, collections *services.CollectionService, log *slog.Logger) *PublicHandler {
	return &PublicHandler{stores: stores, products: products, orders: orders, collections: collections, log: log}
}

func (h *PublicHandler) Store(w http.ResponseWriter, r *http.Request) {
	store, err := h.stores.Public(r.Context(), r.PathValue("slug"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, store)
}

func (h *PublicHandler) Products(w http.ResponseWriter, r *http.Request) {
	products, err := h.products.PublicList(r.Context(), r.PathValue("slug"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, map[string]any{"items": products})
}

func (h *PublicHandler) Product(w http.ResponseWriter, r *http.Request) {
	product, store, err := h.products.PublicGet(r.Context(), r.PathValue("slug"), r.PathValue("id"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, map[string]any{"product": product, "store": store})
}

func (h *PublicHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Customer models.Customer `json:"customer"`
		Items    []struct {
			ProductID      string `json:"productId"`
			VariantID      string `json:"variantId"`
			Quantity       int    `json:"quantity"`
			UnitPriceCents int64  `json:"unitPriceCents"`
		} `json:"items"`
	}
	if err := decode(w, r, &body); err != nil {
		fail(h.log, w, err)
		return
	}
	items := make([]services.ItemInput, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, services.ItemInput{ProductID: item.ProductID, VariantID: item.VariantID, Quantity: item.Quantity, UnitPriceCents: item.UnitPriceCents})
	}
	order, err := h.orders.Place(r.Context(), r.PathValue("slug"), services.PlaceOrderInput{Customer: body.Customer, Items: items})
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusCreated, order)
}

func (h *PublicHandler) Collections(w http.ResponseWriter, r *http.Request) {
	items, err := h.collections.PublicList(r.Context(), r.PathValue("slug"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, map[string]any{"items": items})
}

func (h *PublicHandler) Collection(w http.ResponseWriter, r *http.Request) {
	collection, products, store, err := h.collections.PublicGet(r.Context(), r.PathValue("slug"), r.PathValue("collectionSlug"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, map[string]any{"collection": collection, "products": products, "store": store})
}
