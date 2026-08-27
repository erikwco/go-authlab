package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is a thin wrapper around a pgx connection pool.
type Pool struct {
	inner *pgxpool.Pool
}

// Connect opens a pool and pings Postgres. ctx bounds the ping.
func Connect(ctx context.Context, databaseURL string) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Pool{inner: pool}, nil
}

// Ping reports whether the database is reachable.
func (p *Pool) Ping(ctx context.Context) error {
	if p == nil || p.inner == nil {
		return fmt.Errorf("db pool is nil")
	}
	return p.inner.Ping(ctx)
}

// Close releases pool resources.
func (p *Pool) Close() {
	if p == nil || p.inner == nil {
		return
	}
	p.inner.Close()
}
