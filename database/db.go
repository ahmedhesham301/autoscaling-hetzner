package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	databaseHost, exists := os.LookupEnv("DATABASE_HOST")
	if !exists {
		log.Panic("DATABASE_HOST is not set; expected the Postgres hostname")
	}
	config, err := pgxpool.ParseConfig("postgres://postgres:1234@" + databaseHost + ":5432/postgres")
	if err != nil {
		log.Panicf("failed to parse Postgres connection config for host %q: %v", databaseHost, err)
	}
	config.ConnConfig.TLSConfig = nil

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Panicf("failed to initialize Postgres connection pool for host %q: %v", databaseHost, err)
	}

	_, err = Pool.Exec(context.TODO(), "SELECT now()")
	if err != nil {
		log.Panicf("database connectivity check failed for host %q: %v", databaseHost, err)
	}
}
