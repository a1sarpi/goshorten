package memory_test

import (
	"testing"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage/memory"
)

func TestMemoryStorage_Save(t *testing.T) {
	store := memory.NewMemoryStorage()
	defer store.Close()

	url := &models.URLModel{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
	}

	// Тест сохранения URL
	err := store.Save(url, 1*time.Hour)
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Тест на уникальность оригинального URL
	err = store.Save(url, 1*time.Hour)
	if err != nil {
		t.Errorf("Save() error = %v, expected nil for duplicate URL", err)
	}

	// Тест на сохранение с TTL
	url2 := &models.URLModel{
		ShortCode:   "def456",
		OriginalURL: "https://example2.com",
	}
	err = store.Save(url2, 1*time.Hour)
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}
}

func TestMemoryStorage_Get(t *testing.T) {
	store := memory.NewMemoryStorage()
	defer store.Close()

	url := &models.URLModel{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
	}

	// Сохраняем URL
	err := store.Save(url, 1*time.Hour)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Тест получения существующего URL
	got, err := store.Get("abc123")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if got == nil || got.OriginalURL != url.OriginalURL {
		t.Errorf("Get() = %v, want %v", got, url)
	}

	// Тест получения несуществующего URL
	got, err = store.Get("nonexistent")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if got != nil {
		t.Errorf("Get() = %v, want nil", got)
	}

	// Тест на истечение TTL
	url2 := &models.URLModel{
		ShortCode:   "def456",
		OriginalURL: "https://example2.com",
	}
	err = store.Save(url2, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	got, err = store.Get("def456")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if got != nil {
		t.Errorf("Get() = %v, want nil for expired URL", got)
	}
}

func TestMemoryStorage(t *testing.T) {
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
}
