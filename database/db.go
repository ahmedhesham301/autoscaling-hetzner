package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	databaseHost := os.Getenv("DATABASE_HOST")
	config, err := pgxpool.ParseConfig("postgres://postgres:1234@" + databaseHost + ":5432/postgres")
	if err != nil {
		panic("could not parse connection string: " + err.Error())
	}
	config.ConnConfig.TLSConfig = nil

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic("couldn't connect to database")
	}
}
