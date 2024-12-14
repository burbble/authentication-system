package postgres

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
    pool *pgxpool.Pool
}

func NewPostgres(host, port, user, password, dbname, sslmode string) (*Postgres, error) {
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        user, password, host, port, dbname, sslmode,
    )
    
    config, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }

    return &Postgres{
        pool: pool,
    }, nil
}

func (p *Postgres) Close() error {
    p.pool.Close()
    return nil
}

func (p *Postgres) GetDB() *pgxpool.Pool {
	return p.pool
}
