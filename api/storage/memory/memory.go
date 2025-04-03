package memory

import (
	"encoding/json"
	"errors"
	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage"
	"github.com/a1sarpi/goshorten/pkg/customErrors"
	"os"
	"sync"
	"time"
)

var (
	ErrDuplicatedShortCode = errors.New("short code already exists")
	ErrEmptyShortCode      = errors.New("short code is empty")
	ErrLoadingFile         = errors.New("failed to load file")
)

type Local struct {
	mu   sync.RWMutex
	data map[string]*models.URL
	file string
}

func NewMemoryStorage() *Local {
	return &Local{
		data: make(map[string]*models.URL),
	}
}

func NewWithPersistence(filename string) *Local {
	store := &Local{
		data: make(map[string]*models.URL),
		file: filename,
	}
	store.load()
	return store
}

func (ls *Local) Save(url *models.URL) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if _, exists := ls.data[url.ShortCode]; exists {
		return ErrDuplicatedShortCode
	}

	url.CreatedAt = time.Now()
	ls.data[url.ShortCode] = url
	return nil
}

func (ls *Local) Get(code string) (*models.URL, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	url, exists := ls.data[code]
	if !exists {
		return nil, customErrors.ErrURLNotFound
	}
	return url, nil
}

func (ls *Local) load() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	data, err := os.ReadFile(ls.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &ls.data)
}

func (ls *Local) save() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	data, err := json.MarshalIndent(ls.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ls.file, data, 0644)
}

func (ls *Local) Close() error {
	return ls.save()
}

var _ storage.Storage = (*Local)(nil)
