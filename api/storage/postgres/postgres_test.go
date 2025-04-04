package postgres_test

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage/postgres"
)

// loadEnv загружает переменные окружения из файла .env
func loadEnv() error {
	envFile := filepath.Join(".", ".env")
	file, err := os.Open(envFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func TestPostgresStorage(t *testing.T) {
	// Загружаем переменные окружения из .env
	if err := loadEnv(); err != nil {
		t.Logf("Failed to load .env file: %v", err)
	}

	// Получаем строку подключения из переменной окружения
	connStr := os.Getenv("POSTGRES_CONN_STRING")
	if connStr == "" {
		t.Skip("Skipping PostgreSQL tests: POSTGRES_CONN_STRING not set")
	}

	store, err := postgres.NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("Failed to create PostgreSQL storage: %v", err)
	}
	defer store.Close()

	// Test Save
	url := &models.URLModel{
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
	}
	err = store.Save(url, time.Hour)
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
