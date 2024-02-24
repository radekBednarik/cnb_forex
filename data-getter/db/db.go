package db

// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#pkg-overview
import (
	"context"
	"fmt"
	"log"
	"time"

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

func (d Database) connect() *pgpool.Conn {
	c, err := d.Pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Could not acquire connection from db pool. Error: %v\n", err)
	}

	return c
}

func (dbs Database) CreateTables() {
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

func (dbs Database) insertIntoCountry(value string) string {
	conn := dbs.connect()

	qString := "INSERT INTO country (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = $1 RETURNING id;"

	row := conn.QueryRow(context.Background(), qString, value)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to country table failed with error: %v\n", err)
	}

	return id
}

func (dbs Database) insertIntoCurrName(value string) string {
	conn := dbs.connect()

	qString := "INSERT INTO curr_name (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = $1 RETURNING id;"

	row := conn.QueryRow(context.Background(), qString, value)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to curr_name table failed with error: %v\n", err)
	}

	return id
}

func (dbs Database) insertIntoCurrSymbol(value string) string {
	conn := dbs.connect()

	qString := "INSERT INTO curr_symbol (symbol) VALUES ($1) ON CONFLICT (symbol) DO UPDATE SET symbol = $1 RETURNING id;"

	row := conn.QueryRow(context.Background(), qString, value)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to curr_symbol table failed with error: %v\n", err)
	}

	return id
}

func (dbs Database) insertIntoDate(value string) string {
	conn := dbs.connect()

	qString := "INSERT INTO date (date) VALUES (TO_DATE($1, 'DD.MM.YYYY')) ON CONFLICT (date) DO NOTHING RETURNING id;"

	fmt.Printf("want to insert date: %s\n", value)

	row := conn.QueryRow(context.Background(), qString, value)
	var id string
	err := row.Scan(&id)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to date table failed with error: %v\n", err)
	}

	return id
}

func (dbs Database) SelectIdFromTable(value string, fieldName string, table string) (string, error) {
	conn := dbs.connect()

	// transform the value format to that which db uses for storing DATE
	tParsed, err := time.Parse("01.01.2006", value)
	if err != nil {
		log.Fatalf("Parsing %s value to time failed with error: %v\n", value, err)
	}
	fValue := tParsed.Format("2024-01-01")

	qString := fmt.Sprintf("SELECT id FROM %s WHERE %s = $1 LIMIT 1;", table, fieldName)

	res := conn.QueryRow(context.Background(), qString, fValue)

	var result string
	err = res.Scan(&result)

	conn.Release()

	return result, err
}

func (dbs Database) insertIntoData(countryIndex string, currNameIndex string, currSymbolIndex string, dateIndex string, value float64) {
	conn := dbs.connect()

	qString := "INSERT INTO data (country_id, curr_name_id, curr_symbol_id, date_id, value) VALUES ($1, $2, $3, $4, $5);"

	_, err := conn.Exec(context.Background(), qString, countryIndex, currNameIndex, currSymbolIndex, dateIndex, value)
	conn.Release()

	if err != nil {
		log.Fatalf("Inserting data to 'data' table failed with error: %v\n", err)
	}
}

func (dbs Database) ProcessDailyData(data *p.ForexDataForDate) {
	// check if date from data is already in db table 'date'
	_, err := dbs.SelectIdFromTable(data.Date, "date", "date")
	// if id was found, then data should be already in db and we can exit
	if err == nil {
		return
	}

	// data are not in the db, so do insertions
	idDate := dbs.insertIntoDate(data.Date)

	for _, curr := range data.ForexData {
		idCountry := dbs.insertIntoCountry(curr.Country)
		idCurrName := dbs.insertIntoCurrName(curr.Name)
		idCurrSymbol := dbs.insertIntoCurrSymbol(curr.Symbol)
		dbs.insertIntoData(idCountry, idCurrName, idCurrSymbol, idDate, curr.Value)
	}
}
