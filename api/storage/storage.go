package storage

import (
	"github.com/a1sarpi/goshorten/api/models"
)

type Storage interface {
	Save(url *models.URL) error
	Get(shortCode string) (*models.URL, error)
}
