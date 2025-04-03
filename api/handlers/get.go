package handlers

import (
	"github.com/a1sarpi/goshorten/api/storage"
	"net/http"
)

type GetHandler struct {
	storage storage.Storage
}

func NewGetHandler(storage storage.Storage) *GetHandler {
	return &GetHandler{storage: storage}
}

func (h *GetHandler) HandleRedirect(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, "[ERROR] Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.Error(rw, "[ERROR] Short code not provided", http.StatusBadRequest)
		return
	}

	url, err := h.storage.Get(shortCode)
	if err != nil {
		http.Error(rw, "[ERROR] URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(rw, r, url.OriginalURL, http.StatusMovedPermanently)
}
