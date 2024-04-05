package database

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	Pool *pgxpool.Pool
	once sync.Once
}

func (p *PostgresClient) NewPostgresClient(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	var err error
	p.once.Do(func() {
		p.Pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return p.Pool, err
}

func (p *PostgresClient) PingDB(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}
