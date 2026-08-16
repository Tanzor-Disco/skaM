package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func Connect(URI string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), URI)
	if err != nil {
		return &DB{}, err
	}
	return &DB{
		pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
