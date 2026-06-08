package core

import (
	"context"
	"fmt"
	"time"

	"github.com/vpa/quanlynhahang-backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolCfg.MaxConns = 6
	poolCfg.MinConns = 0
	poolCfg.MaxConnLifetime = 30 * time.Minute
	return pgxpool.NewWithConfig(ctx, poolCfg)
}
