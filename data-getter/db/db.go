package db

// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#pkg-overview
import (
	"context"
	"fmt"
	"log"

	pgpool "github.com/jackc/pgx/v5/pgxpool"
	p "github.com/radekBednarik/cnb_forex/data-getter/parser"
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
	conn.Release()

	if err != nil {
		log.Fatalf("Scanning query result for value failed with error: %v", err)
	}

	// if the db is empty, create tables
	if !result {
		// ensure, that primary uuid keys are always automatically created by default
		statement := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
		conn = dbs.connect()
		_, err := conn.Exec(context.Background(), statement)
		conn.Release()

		if err != nil {
			log.Fatalf("Creating uuid-ossp extension failed with error: %v", err)
		}

		statements := []string{
			`-- Creating the country table
CREATE TABLE country (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL
);
`,
			`-- Creating the curr_name table
CREATE TABLE curr_name (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL
);
`,
			`-- Creating the curr_symbol table
CREATE TABLE curr_symbol (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    symbol VARCHAR UNIQUE NOT NULL
);
`,
			`-- Creating the date table
CREATE TABLE date (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    date DATE UNIQUE NOT NULL
);
`,
			`-- Creating the data table with foreign key constraints
CREATE TABLE data (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    country_id UUID REFERENCES country(id),
    curr_name_id UUID REFERENCES curr_name(id),
    curr_symbol_id UUID REFERENCES curr_symbol(id),
    date_id UUID REFERENCES date(id),
    value FLOAT NOT NULL
);
`,
		}

		for _, statement := range statements {
			conn := dbs.connect()
			_, err := conn.Exec(context.Background(), statement)
			conn.Release()

			if err != nil {
				log.Fatalf("Creating tables in db failed with error: %v\n", err)
			}
		}

	}
}

func InsertIntoCountry(value string, dbs Database) string {
	conn := dbs.connect()

	qString := fmt.Sprintf("INSERT INTO country (name) VALUES (%s) ON CONFLICT (name) DO NOTHING RETURNING id;", value)

	row := conn.QueryRow(context.Background(), qString)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to country table failed with error: %v\n", err)
	}

	return id
}

func InsertIntoCurrName(value string, dbs Database) string {
	conn := dbs.connect()

	qString := fmt.Sprintf("INSERT INTO curr_name (name) VALUES (%s) ON CONFLICT (name) DO NOTHING RETURNING id;", value)

	row := conn.QueryRow(context.Background(), qString)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to curr_name table failed with error: %v\n", err)
	}

	return id
}

func InsertIntoCurrSymbol(value string, dbs Database) string {
	conn := dbs.connect()

	qString := fmt.Sprintf("INSERT INTO curr_symbol (symbol) VALUES (%s) ON CONFLICT (symbol) DO NOTHING RETURNING id;", value)

	row := conn.QueryRow(context.Background(), qString)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to curr_symbol table failed with error: %v\n", err)
	}

	return id
}

func InsertIntoDate(value string, dbs Database) string {
	conn := dbs.connect()

	qString := fmt.Sprintf("INSERT INTO date (date) VALUES (%s) ON CONFLICT (date) DO NOTHING;", value)

	row := conn.QueryRow(context.Background(), qString)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to date table failed with error: %v\n", err)
	}

	return id
}

func SelectIdFromTable(value string, fieldName string, table string, dbs Database) (string, error) {
	conn := dbs.connect()

	qString := fmt.Sprintf("SELECT id from %s WHERE %s = '%s' LIMIT 1;", table, fieldName, value)

	res := conn.QueryRow(context.Background(), qString)

	var result string
	err := res.Scan(&result)

	conn.Release()

	return result, err
}

func InsertIntoData(countryIndex string, currNameIndex string, currSymbolIndex string, dateIndex string, value float64, dbs Database) {
	conn := dbs.connect()

	qString := fmt.Sprintf("INSERT INTO data (country, curr_name, curr_symbol, date, value) VALUES (%s, %s, %s, %s, %f);", countryIndex, currNameIndex, currSymbolIndex, dateIndex, value)

	_, err := conn.Exec(context.Background(), qString)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to 'data' table failed with error: %v\n", err)
	}
}

func ProcessDailyData(data *p.ForexDataForDate, dbs Database) {
	// check if date from data is already in db table 'date'
	id, err := SelectIdFromTable(data.Date, "date", "date", dbs)
	if err != nil {
		log.Fatalf("Error when processing the daily forex data. Error: %v\n", err)
	}

	// if id was found, then data should be already in db and we can exit
	// this is a very weak check, TODO: do it better
	if id != "" {
		return
	}

	// data are not in the db, so do insertions
	idDate := InsertIntoDate(data.Date, dbs)

	for _, curr := range data.ForexData {
		idCountry := InsertIntoCountry(curr.Country, dbs)
		idCurrName := InsertIntoCurrName(curr.Name, dbs)
		idCurrSymbol := InsertIntoCurrSymbol(curr.Symbol, dbs)

	}
}
