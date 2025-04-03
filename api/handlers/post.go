package handlers

import (
	"encoding/json"
	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	allowedLength  = 10
	allowedLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var (
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type PostHandler struct {
	storage storage.Storage
}

func NewPostHandler(storage storage.Storage) *PostHandler {
	return &PostHandler{storage: storage}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (ph *PostHandler) HandleShorten(rw http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request: Method=%s, Path=%s", r.Method, r.URL.Path)

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST")

	if r.Method == http.MethodOptions {
		rw.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("Error: Method %s not allowed", r.Method)
		http.Error(rw, "[ERROR] Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(rw, "[ERROR] Invalid JSON", http.StatusBadRequest)
		return
	}

	shortCode, err := ph.shortenURL(req.URL)
	if err != nil {
		log.Printf("Error creating short URL: %v", err)
		http.Error(rw, "[ERROR] Failed to create short URL", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(map[string]string{
		"short_url": shortCode,
	})
}

func (ph *PostHandler) shortenURL(originalURL string) (string, error) {
	shortCode := generateShortCode()
	url := &models.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}

	if err := ph.storage.Save(url); err != nil {
		return "", err
	}

	return shortCode, nil
}

func generateShortCode() string {
	code := make([]byte, allowedLength)
	for i := range code {
		code[i] = allowedLetters[randomGenerator.Intn(len(allowedLetters))]
	}

	return string(code)
}
