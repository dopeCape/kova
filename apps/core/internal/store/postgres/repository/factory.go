package repository

import (
	"context"
	"log"

	"github.com/dopeCape/kova/internal/config"
	"github.com/dopeCape/kova/internal/shared/database"
	"github.com/dopeCape/kova/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDefaultStore(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, store.Store) {
	dbConfig := &database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	store := NewStore(db)
	return db, store
}
