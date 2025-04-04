package storage

import (
	"time"

	"github.com/a1sarpi/goshorten/api/models"
)

type Storage interface {
	Save(url *models.URLModel, ttl time.Duration) error
	Get(shortCode string) (*models.URLModel, error)
	Close() error
}

// NewPostgresStorage is a factory function that creates a new PostgreSQL storage
// This is just a declaration, the actual implementation is in the postgres package
func NewPostgresStorage(connStr string) (Storage, error) {
	return nil, nil
}
