package repository

import (
	"context"

	"github.com/dopeCape/kova/internal/store"
	"github.com/dopeCape/kova/internal/store/postgres/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

func NewStore(db *pgxpool.Pool) store.Store {
	return &Store{
		db:      db,
		queries: generated.New(db),
	}
}

func (s *Store) WithTx(ctx context.Context, fn func(store.Store) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txStore := &Store{
		db:      s.db,
		queries: s.queries.WithTx(tx),
	}

	if err := fn(txStore); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *Store) Close() {
	s.db.Close()
}
