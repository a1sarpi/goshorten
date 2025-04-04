package storage_test

import (
	"testing"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage/memory"
)

// TestStorage is a helper function to test storage implementations
func TestStorage(t *testing.T) {
	t.Run("Storage Tests", func(t *testing.T) {
		store := memory.NewMemoryStorage()
		defer store.Close()

		// Test Save
		url := &models.URLModel{
			ShortCode:   "test123",
			OriginalURL: "https://example.com",
		}
		err := store.Save(url, time.Hour)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		// Test Get
		got, err := store.Get("test123")
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if got == nil || got.OriginalURL != url.OriginalURL {
			t.Errorf("Get() = %v, want %v", got, url)
		}

		// Test Get non-existent
		got, err = store.Get("nonexistent")
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if got != nil {
			t.Errorf("Get() = %v, want nil", got)
		}

		// Test TTL
		url2 := &models.URLModel{
			ShortCode:   "test456",
			OriginalURL: "https://example2.com",
		}
		err = store.Save(url2, time.Millisecond)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}
		time.Sleep(2 * time.Millisecond)
		got, err = store.Get("test456")
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if got != nil {
			t.Errorf("Get() = %v, want nil for expired URL", got)
		}
	})
}
