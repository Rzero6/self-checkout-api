package config

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

const DBTimeout = 10 * time.Second

func ConnectDB() {
	ctx, cancel := DBContext()
	defer cancel()

	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	// Pool config
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	DB, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// Test connection
	if err := DB.Ping(ctx); err != nil {
		panic(err)
	}
}

func DBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DBTimeout)
}
