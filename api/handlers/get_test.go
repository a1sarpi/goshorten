package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		shortCode      string
		url            *models.URLModel
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Valid shortcode",
			shortCode: "abc123",
			url: &models.URLModel{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				CreatedAt:   time.Now().Unix(),
				ExpiresAt:   time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:           "Empty shortcode",
			shortCode:      "",
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.ErrorResponse{
				Message: "Shortcode is required",
				Code:    http.StatusBadRequest,
			},
		},
		{
			name:           "Non-existent shortcode",
			shortCode:      "nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody: models.ErrorResponse{
				Message: "URL not found",
				Code:    http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			store := NewMockStorage()
			if tt.url != nil {
				store.Save(tt.url, 24*time.Hour)
			}
			handler := NewGetHandler(store)

			// Create test router
			router := gin.New()
			router.GET("/:shortcode", handler.HandleRedirect)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/"+tt.shortCode, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body for error cases
			if tt.expectedStatus != http.StatusMovedPermanently {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(models.ErrorResponse).Message, response.Message)
				assert.Equal(t, tt.expectedBody.(models.ErrorResponse).Code, response.Code)
			} else {
				// Check redirect location
				assert.Equal(t, tt.url.OriginalURL, w.Header().Get("Location"))
			}
		})
	}
}
