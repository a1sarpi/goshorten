package main

import (
	"fmt"
	"log"
	"os"

	"github.com/a1sarpi/goshorten/api"
	"github.com/a1sarpi/goshorten/api/storage"
	"github.com/a1sarpi/goshorten/api/storage/memory"
	"github.com/a1sarpi/goshorten/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage
	var store storage.Storage
	if cfg.Storage.Type == "postgres" {
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
			cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
		store, err = storage.NewPostgresStorage(connStr)
		if err != nil {
			log.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
		}
	} else {
		store = memory.NewMemoryStorage()
	}
	defer store.Close()

	// Create and start server
	server := api.NewServer(store)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := server.Start(addr); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
