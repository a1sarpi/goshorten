package handlers

import (
	"net/http"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage"
	"github.com/gin-gonic/gin"
)

type GetHandler struct {
	storage storage.Storage
}

func NewGetHandler(storage storage.Storage) *GetHandler {
	return &GetHandler{
		storage: storage,
	}
}

// @Summary Redirect to original URL
// @Description Redirects to the original URL using the short code. Returns 404 if the URL is not found or has expired.
// @Tags urls
// @Param shortcode path string true "Short URL code" example(abc123)
// @Success 301 {string} string "Redirect to original URL"
// @Failure 400 {object} models.ErrorResponse "Invalid request - empty shortcode"
// @Failure 404 {object} models.ErrorResponse "URL not found or expired"
// @Failure 410 {object} models.ErrorResponse "URL has expired"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /{shortcode} [get]
func (gh *GetHandler) HandleRedirect(c *gin.Context) {
	shortCode := c.Param("shortcode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Shortcode is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	url, err := gh.storage.Get(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to fetch URL",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if url == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "URL not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}
