package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	databaseHost, exists := os.LookupEnv("DATABASE_HOST")
	if !exists {
		slog.Error("env var DATABASE_HOST is not set")
		os.Exit(1)
	}

	config, err := pgxpool.ParseConfig("postgres://postgres:1234@" + databaseHost + ":5432/postgres")
	if err != nil {
		slog.Error("failed to parse Postgres connection config", "error", err)
		os.Exit(1)
	}
	config.ConnConfig.TLSConfig = nil

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("failed to initialize Postgres connection pool", "error", err)
		os.Exit(1)
	}

	_, err = Pool.Exec(context.TODO(), "SELECT now()")
	if err != nil {
		slog.Error("database connectivity check failed for host", "error", err)
		os.Exit(1)
	}
}
