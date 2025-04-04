package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockStorage реализует интерфейс storage.Storage для тестов
type MockStorage struct {
	urls map[string]*models.URLModel
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		urls: make(map[string]*models.URLModel),
	}
}

func (ms *MockStorage) Save(url *models.URLModel, ttl time.Duration) error {
	url.CreatedAt = time.Now().Unix()
	url.ExpiresAt = time.Now().Add(ttl).Unix()
	ms.urls[url.ShortCode] = url
	return nil
}

func (ms *MockStorage) Get(shortCode string) (*models.URLModel, error) {
	if url, ok := ms.urls[shortCode]; ok {
		if time.Now().Unix() > url.ExpiresAt {
			delete(ms.urls, shortCode)
			return nil, nil
		}
		return url, nil
	}
	return nil, nil
}

func (ms *MockStorage) Close() error {
	return nil
}

func TestPostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Valid URL",
			requestBody: models.URLRequest{
				URL: "https://example.com",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: models.URLResponse{
				OriginalURL: "https://example.com",
			},
		},
		{
			name: "Empty URL",
			requestBody: models.URLRequest{
				URL: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.ErrorResponse{
				Message: "Invalid request body",
				Code:    http.StatusBadRequest,
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody: models.ErrorResponse{
				Message: "Invalid request body",
				Code:    http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			store := NewMockStorage()
			handler := NewPostHandler(store)

			// Create test router
			router := gin.New()
			router.POST("/shorten", handler.HandleShorten)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body
			if tt.expectedStatus == http.StatusCreated {
				var response models.URLResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(models.URLResponse).OriginalURL, response.OriginalURL)
				assert.NotEmpty(t, response.ShortURL)
				assert.NotZero(t, response.ExpiresAt)
			} else {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(models.ErrorResponse).Message, response.Message)
				assert.Equal(t, tt.expectedBody.(models.ErrorResponse).Code, response.Code)
			}
		})
	}
}
