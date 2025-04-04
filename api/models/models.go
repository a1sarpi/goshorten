package models

// @Description Request body for URL shortening
type URLRequest struct {
	// @required
	// @example https://example.com
	URL string `json:"url" binding:"required,url"`

	// @example 3600
	TTL int64 `json:"ttl,omitempty"`
}

// @Description Response for URL shortening
type URLResponse struct {
	// @example http://localhost:8080/abc123
	ShortURL string `json:"short_url"`

	// @example https://example.com
	OriginalURL string `json:"original_url"`

	// @example 1672531200
	ExpiresAt int64 `json:"expires_at"`
}

// @Description Error response
type ErrorResponse struct {
	// @example Invalid URL format
	Message string `json:"message"`

	// @example 400
	Code int `json:"code"`
}

type URLModel struct {
	ShortCode   string
	OriginalURL string
	CreatedAt   int64
	ExpiresAt   int64
}
