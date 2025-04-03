package postgres

//
//import (
//	"errors"
//	"github.com/a1sarpi/goshorten/api/models"
//	"github.com/a1sarpi/goshorten/api/storage"
//	"github.com/a1sarpi/goshorten/pkg/customErrors"
//	"gorm.io/gorm"
//	"time"
//)
//
//var _ storage.Storage = (*SQLStorage)(nil)
//
//type SQLStorage struct {
//	db *gorm.DB
//}
//
//func NewSQLStorage(db *gorm.DB) (*SQLStorage, error) {
//	if err := db.AutoMigrate(&SQLStorage{}); err != nil {
//		return nil, err
//	}
//
//	return &SQLStorage{db: db}, nil
//}
//
//func (s *SQLStorage) Save(url *models.URL) error {
//	result := s.db.Create(url)
//	return result.Error
//}
//
//func (s *SQLStorage) FindByShortCode(code string) (*models.URL, error) {
//	var url models.URL
//	result := s.db.Where("short_code = ?", code).First(&url)
//
//	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//		return nil, customErrors.ErrURLNotFound
//	}
//
//	if url.ExpiresAt.Before(time.Now()) {
//		s.db.Delete(&url)
//		return nil, customErrors.ErrURLExpired
//	}
//
//	return &url, nil
//}
//
//func (s *SQLStorage) FindByOriginalURL(originalURL string) (*models.URL, error) {
//	var url models.URL
//	result := s.db.Where("original_url = ?", originalURL).First(&url)
//
//	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//		return nil, customErrors.ErrURLNotFound
//	}
//
//	return &url, result.Error
//}
//
//func (s *SQLStorage) Delete(code *models.URL) error {
//	return s.db.Where("short_code = ?", code).Delete(&models.URL{}).Error
//}
//
//func (s *SQLStorage) PurgeExpiredURLs() (int64, error) {
//	result := s.db.Where("expires_at < ?", time.Now()).Delete(&models.URL{})
//	return result.RowsAffected, result.Error
//}
//
//func (s *SQLStorage) IncrementClicks(code string) error {
//	return s.db.Model(&models.URL{}).
//		Where("short_code = ?", code).
//		Update("clicks", gorm.Expr("clicks + 1")).
//		Error
//}
