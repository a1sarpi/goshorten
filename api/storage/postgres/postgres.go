package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/a1sarpi/goshorten/api/models"
	"github.com/a1sarpi/goshorten/api/storage"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (storage.Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (ps *PostgresStorage) Save(url *models.URLModel, ttl time.Duration) error {
	ctx := context.Background()

	var existingShortCode string
	err := ps.db.QueryRowContext(ctx,
		"SELECT short_code FROM urls WHERE original_url = $1",
		url.OriginalURL).Scan(&existingShortCode)

	if err == nil {
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	_, err = ps.db.ExecContext(ctx,
		"INSERT INTO urls (short_code, original_url, expires_at) VALUES ($1, $2, $3)",
		url.ShortCode,
		url.OriginalURL,
		time.Now().Add(ttl))

	return err
}

func (ps *PostgresStorage) Get(shortCode string) (*models.URLModel, error) {
	ctx := context.Background()

	var url models.URLModel
	var expiresAt sql.NullTime

	err := ps.db.QueryRowContext(ctx,
		"SELECT short_code, original_url, created_at, expires_at FROM urls WHERE short_code = $1",
		shortCode).Scan(&url.ShortCode, &url.OriginalURL, &url.CreatedAt, &expiresAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid && time.Now().After(expiresAt.Time) {
		return nil, nil
	}

	return &url, nil
}

func (ps *PostgresStorage) Close() error {
	return ps.db.Close()
}
