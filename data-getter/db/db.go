package db

// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#pkg-overview
import (
	"context"
	"fmt"
	"log"

	pgpool "github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	connString string
	connConfig *pgpool.Config
	pool       *pgpool.Pool
}

func (d *Db) New(connString string) *Db {

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
	d.pool = p

	return d

}

func (d *Db) connect() *pgpool.Conn {
	c, err := d.pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Could not acquire connection from db pool. Error: %v\n", err)
	}

	return c
}

func CreateTables(dbPool Db) {
	conn := dbPool.connect()

	qString := `
		SELECT EXISTS (
   			SELECT 1
   			FROM information_schema.tables
   			WHERE table_schema = 'public'
   			LIMIT 1
		);
	`

	res := conn.QueryRow(context.Background(), qString)

	fmt.Println(res)
}
