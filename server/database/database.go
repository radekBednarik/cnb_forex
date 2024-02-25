package database

import (
	"context"
	"log"

	pgpool "github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	connConfig *pgpool.Config
	Pool       *pgpool.Pool
	connString string
}

func (d *Database) New(connString string) *Database {
	d.connString = connString
	c, err := pgpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Failed to parse db config with error: %v\n", err)
	}
	d.connConfig = c

	p, err := pgpool.NewWithConfig(context.Background(), d.connConfig)
	if err != nil {
		log.Fatalf("Failed to create new pool connection to db with error: %v\n", err)
	}
	d.Pool = p

	return d
}

func (d Database) connect() *pgpool.Conn {
	c, err := d.Pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Could not acquire connection from db pool. Error: %v\n", err)
	}

	return c
}
