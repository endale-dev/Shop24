package handlers

import (
	"log/slog"
	"net/http"

	"shop24/internal/httpx"
	"shop24/internal/services"
)

type CollectionHandler struct {
	collections *services.CollectionService
	log         *slog.Logger
}

func NewCollectionHandler(collections *services.CollectionService, log *slog.Logger) *CollectionHandler {
	return &CollectionHandler{collections: collections, log: log}
}

func (h *CollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	body, err := collectionBody(w, r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	collection, err := h.collections.Create(r.Context(), owner, body)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusCreated, collection)
}

func (h *CollectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	body, err := collectionBody(w, r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	collection, err := h.collections.Update(r.Context(), owner, r.PathValue("id"), body)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, collection)
}

func (h *CollectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	if err := h.collections.Delete(r.Context(), owner, r.PathValue("id")); err != nil {
		fail(h.log, w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionHandler) Get(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	collection, err := h.collections.Get(r.Context(), owner, r.PathValue("id"))
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, collection)
}

func (h *CollectionHandler) List(w http.ResponseWriter, r *http.Request) {
	owner, err := currentUser(r)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	items, err := h.collections.List(r.Context(), owner)
	if err != nil {
		fail(h.log, w, err)
		return
	}
	httpx.Data(w, http.StatusOK, map[string]any{"items": items})
}

func collectionBody(w http.ResponseWriter, r *http.Request) (services.CollectionInput, error) {
	var body struct {
		Name        string   `json:"name"`
		Slug        string   `json:"slug"`
		Description string   `json:"description"`
		ImageURL    string   `json:"imageUrl"`
		ProductIDs  []string `json:"productIds"`
		Published   bool     `json:"published"`
	}
	if err := decode(w, r, &body); err != nil {
		return services.CollectionInput{}, err
	}
	return services.CollectionInput{
		Name: body.Name, Slug: body.Slug, Description: body.Description,
		ImageURL: body.ImageURL, ProductIDs: body.ProductIDs, Published: body.Published,
	}, nil
}
