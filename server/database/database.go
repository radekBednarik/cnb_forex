package database

import (
	"context"
	"log"
	"time"

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

type SingleCurrData struct {
	Name    string
	Symbol  string
	Country string
	Value   float64
}

type DataByDate = map[string][]SingleCurrData

type Data struct {
	Data DataByDate
}

func (d Database) SelectDashboardDataV1(dateFrom string, dateTo string) (Data, error) {
	conn := d.connect()
	defer conn.Release()

	statement := `
    select 
      dt.date as date,
      c."name" as country_name,
      cn."name" as currency_name,
      cs.symbol as currency_symbol,
      d.value as czk_to_currency_value
    from "data" d 
    left join country c 
    on d.country_id = c.id 
    left join curr_name cn 
    on d.curr_name_id = cn.id 
    left join curr_symbol cs 
    on d.curr_symbol_id = cs.id
    left join "date" dt
    on d.date_id = dt.id
    where dt.date between $1 and $2;
  `

	rows, err := conn.Query(context.Background(), statement, dateFrom, dateTo)
	if err != nil {
		return Data{}, err
	}
	defer rows.Close()

	data := Data{}
	dataByDate := DataByDate{}
	singleDateData := []SingleCurrData{}
	tempDate := dateFrom

	for rows.Next() {
		var date time.Time
		var currData SingleCurrData

		err := rows.Scan(&date, &currData.Country, &currData.Name, &currData.Symbol, &currData.Value)
		if err != nil {
			return Data{}, err
		}

		fDate := date.Format("2006-01-02")

		if tempDate == fDate {
			singleDateData = append(singleDateData, currData)
			continue
		}

		// date changed
		// add data list to map date prop
		dataByDate[fDate] = singleDateData
		// clear list
		singleDateData = singleDateData[:0]
		// fill new value
		singleDateData = append(singleDateData, currData)
		// adjust tempDate to new date
		tempDate = fDate

	}

	data.Data = dataByDate

	return data, nil
}
