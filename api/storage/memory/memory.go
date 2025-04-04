package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage"
)

var (
	ErrDuplicatedShortCode = errors.New("short code already exists")
	ErrEmptyShortCode      = errors.New("short code is empty")
	ErrLoadingFile         = errors.New("failed to load file")
)

type MemoryStorage struct {
	mu    sync.RWMutex
	urls  map[string]*models.URLModel
	ttls  map[string]time.Time
	clean chan struct{}
}

func NewMemoryStorage() *MemoryStorage {
	ms := &MemoryStorage{
		urls:  make(map[string]*models.URLModel),
		ttls:  make(map[string]time.Time),
		clean: make(chan struct{}),
	}
	go ms.cleanupLoop()
	return ms
}

func (ms *MemoryStorage) Save(url *models.URLModel, ttl time.Duration) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, existingURL := range ms.urls {
		if existingURL.OriginalURL == url.OriginalURL {
			return nil
		}
	}

	ms.urls[url.ShortCode] = url
	if ttl > 0 {
		ms.ttls[url.ShortCode] = time.Now().Add(ttl)
	}
	return nil
}

func (ms *MemoryStorage) Get(shortCode string) (*models.URLModel, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	url, exists := ms.urls[shortCode]
	if !exists {
		return nil, nil
	}

	if ttl, hasTTL := ms.ttls[shortCode]; hasTTL {
		if time.Now().After(ttl) {
			return nil, nil
		}
	}

	return url, nil
}

func (ms *MemoryStorage) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ms.cleanup()
		case <-ms.clean:
			return
		}
	}
}

func (ms *MemoryStorage) cleanup() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now()
	for shortCode, ttl := range ms.ttls {
		if now.After(ttl) {
			delete(ms.urls, shortCode)
			delete(ms.ttls, shortCode)
		}
	}
}

func (ms *MemoryStorage) Close() error {
	close(ms.clean)
	return nil
}

var _ storage.Storage = (*MemoryStorage)(nil)
