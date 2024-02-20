package db

// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#pkg-overview
import (
	"context"
	"log"

	pgpool "github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	connString string
	connConfig *pgpool.Config
	Pool       *pgpool.Pool
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

func (d *Database) connect() *pgpool.Conn {
	c, err := d.Pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Could not acquire connection from db pool. Error: %v\n", err)
	}

	return c
}

func CreateTables(dbs Database) {
	conn := dbs.connect()
	// check, if there is the main table in the db
	qString := `
		SELECT EXISTS (
   			SELECT FROM information_schema.tables
   			WHERE table_schema LIKE 'public'
				AND table_type LIKE 'BASE TABLE'
				AND table_name = 'data'
		);
	`

	res := conn.QueryRow(context.Background(), qString)
	var result bool

	err := res.Scan(&result)
	if err != nil {
		log.Fatalf("Scanning query result for value failed with error: %v", err)
	}

	// if the db is empty, create tables
	if !result {

		qString = `
-- Creating the country table
CREATE TABLE country (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL
);

-- Creating the curr_name table
CREATE TABLE curr_name (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL
);

-- Creating the curr_symbol table
CREATE TABLE curr_symbol (
    id UUID PRIMARY KEY,
    symbol VARCHAR NOT NULL
);

-- Creating the date table
CREATE TABLE date (
    id UUID PRIMARY KEY,
    date DATE NOT NULL
);

-- Creating the data table with foreign key constraints
CREATE TABLE pub.data (
    id UUID PRIMARY KEY,
    country_id UUID REFERENCES pub.country(id),
    curr_name_id UUID REFERENCES pub.curr_name(id),
    curr_symbol_id UUID REFERENCES pub.curr_symbol(id),
    date_id UUID REFERENCES pub.date(id),
    value FLOAT NOT NULL
);
`
		_, err := conn.Query(context.Background(), qString)

		if err != nil {
			log.Fatalf("Creating tables in db failed with error: %v\n", err)
		}
	}

}
