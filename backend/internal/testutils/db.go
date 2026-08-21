package testutils

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TestDB struct {
	pool *pgxpool.Pool
}

func Connect(URI string) (*TestDB, error) {
	pool, err := pgxpool.Connect(context.Background(), URI)
	if err != nil {
		return &TestDB{}, err
	}
	return &TestDB{
		pool: pool,
	}, err
}

func (tdb *TestDB) Close() {
	tdb.pool.Close()
}
