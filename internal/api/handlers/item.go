package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"nethttppractice/internal/domain"
)

type ItemHandler struct {
	repo domain.Repository
}

func NewItemHandler(repo domain.Repository) *ItemHandler {
	return &ItemHandler{
		repo: repo,
	}
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	items, err := h.repo.GetAllItems(r.Context())
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) InsertItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var item domain.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "body error", http.StatusInternalServerError)
		slog.Error("decode json", "err", err.Error())
		return
	}

	if err := h.repo.Create(r.Context(), &item); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		slog.Error("db error", "err", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}