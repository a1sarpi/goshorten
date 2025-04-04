package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage"
	"github.com/gin-gonic/gin"
)

// @title GoShorten API
// @version 1.0
// @description URL shortening service with support for in-memory and PostgreSQL storage
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @Summary Create a short URL
// @Description Creates a shortened URL from the provided original URL
// @Tags urls
// @Accept json
// @Produce json
// @Param request body models.URLRequest true "URL to shorten"
// @Success 201 {object} models.URLResponse "URL successfully shortened"
// @Failure 400 {object} models.ErrorResponse "Invalid request - empty URL or invalid format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - invalid or missing API key"
// @Failure 403 {object} models.ErrorResponse "Forbidden - rate limit exceeded"
// @Failure 409 {object} models.ErrorResponse "Conflict - URL already exists"
// @Failure 422 {object} models.ErrorResponse "Unprocessable Entity - invalid TTL value"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /shorten [post]

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
	return &PostHandler{
		storage: storage,
	}
}

func (ph *PostHandler) HandleShorten(c *gin.Context) {
	var req models.URLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Generate short code
	shortCode := generateShortCode()

	// Set default TTL if not provided
	ttl := time.Duration(req.TTL) * time.Second
	if ttl == 0 {
		ttl = 24 * time.Hour // Default TTL
	}

	// Create URL model
	url := &models.URLModel{
		ShortCode:   shortCode,
		OriginalURL: req.URL,
		CreatedAt:   time.Now().Unix(),
		ExpiresAt:   time.Now().Add(ttl).Unix(),
	}

	// Save to storage
	if err := ph.storage.Save(url, ttl); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to save URL",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, models.URLResponse{
		ShortURL:    "http://" + c.Request.Host + "/" + shortCode,
		OriginalURL: req.URL,
		ExpiresAt:   url.ExpiresAt,
	})
}

func generateShortCode() string {
	code := make([]byte, allowedLength)
	for i := range code {
		code[i] = allowedLetters[randomGenerator.Intn(len(allowedLetters))]
	}
	return string(code)
}
