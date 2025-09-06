package apihttp

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/yourname/inventory-go/internal/model"
	"github.com/yourname/inventory-go/internal/store"
	"github.com/yourname/inventory-go/internal/validate"
)

type API struct {
	Store *store.Store
}

func Router(s *store.Store) http.Handler {
	api := &API{Store: s}
	r := chi.NewRouter()

	r.Post("/items", api.handleCreateItem)
	r.Get("/items", api.handleListItems)
	r.Get("/items/{id}", api.handleGetItem)
	r.Patch("/items/{id}", api.handlePatchItem)
	r.Delete("/items/{id}", api.handleDeleteItem)

	return r
}

type errorResp struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (a *API) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var in model.ItemCreate
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorResp{Error: "invalid JSON"})
		return
	}
	if err := validate.ValidateCreateItem(&in); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorResp{Error: "validation error", Details: err})
		return
	}
	it, err := a.Store.CreateItem(r.Context(), in)
	if err != nil {
		if err == store.ErrConflict {
			writeJSON(w, http.StatusConflict, errorResp{Error: "name already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResp{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, it)
}

func (a *API) handleListItems(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limitStr := q.Get("limit")
	offsetStr := q.Get("offset")

	limit := 100 // по умолчанию
	offset := 0

	if offsetStr != "" {
		v, err := strconv.Atoi(offsetStr)
		if err != nil || v < 0 {
			writeJSON(w, http.StatusUnprocessableEntity, errorResp{Error: "invalid offset"})
			return
		}
		offset = v
	}

	if limitStr != "" {
		v, err := strconv.Atoi(limitStr)
		if err != nil || v < 1 || v > 1000 {
			writeJSON(w, http.StatusUnprocessableEntity, errorResp{Error: "invalid limit"})
			return
		}
		limit = v
	}

	items, err := a.Store.ListItems(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResp{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) handleGetItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	it, err := a.Store.GetItemByID(r.Context(), id)
	if err != nil {
		if err == store.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResp{Error: "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResp{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (a *API) handlePatchItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var in model.ItemUpdate
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorResp{Error: "invalid JSON"})
		return
	}
	if err := validate.ValidateUpdateItem(&in); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorResp{Error: "validation error", Details: err})
		return
	}
	it, err := a.Store.UpdateItemPartial(r.Context(), id, in)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			writeJSON(w, http.StatusNotFound, errorResp{Error: "not found"})
		case store.ErrConflict:
			writeJSON(w, http.StatusConflict, errorResp{Error: "name already exists"})
		default:
			writeJSON(w, http.StatusInternalServerError, errorResp{Error: "internal error"})
		}
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func (a *API) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := a.Store.DeleteItem(r.Context(), id); err != nil {
		if err == store.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResp{Error: "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResp{Error: "internal error"})
		return
	}
	// 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusUnprocessableEntity, errorResp{Error: "invalid id"})
		return 0, false
	}
	return id, true
}

// Helper for tests (allows injecting context with timeout if needed)
func withCtx(r *http.Request, ctx context.Context) *http.Request {
	return r.WithContext(ctx)
}
