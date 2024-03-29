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
	Name    string  `json:"name"`
	Symbol  string  `json:"symbol"`
	Country string  `json:"country"`
	Value   float64 `json:"value"`
}

type DataByDate = map[string][]SingleCurrData

type Data struct {
	Data DataByDate `json:"data"`
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
    where dt.date between $1 and $2
    order by dt.date asc;
  `

	rows, err := conn.Query(context.Background(), statement, dateFrom, dateTo)
	if err != nil {
		return Data{}, err
	}
	defer rows.Close()

	data := Data{}
	dataByDate := DataByDate{}
	singleDateData := []SingleCurrData{}
	tempDate := dateTo

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
			// handle special case, when dateFrom and dateTo are equal
			if fDate == dateFrom {
				dataByDate[fDate] = singleDateData
			}

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

	if rows.Err() != nil {
		return Data{}, rows.Err()
	}

	data.Data = dataByDate

	return data, nil
}

type OneCurrData struct {
	Currency string    `json:"currency"`
	Dates    []string  `json:"dates"`
	Values   []float64 `json:"values"`
}

func (d Database) SelectDashboardDataV2(dateFrom string, dateTo string, currency string) (OneCurrData, error) {
	conn := d.connect()
	defer conn.Release()

	statement := `
    select
      dt."date" as date, 
      round(d.value::numeric, 3) as value
    from "data" d 
    left join "date" dt
    on d.date_id = dt.id 
    left join curr_symbol cs 
    on d.curr_symbol_id = cs.id 
    where dt."date" between $1 and $2 
      and cs.symbol = $3 
    order by dt."date" asc;
  `

	rows, err := conn.Query(context.Background(), statement, dateFrom, dateTo, currency)
	if err != nil {
		return OneCurrData{}, err
	}
	defer rows.Close()

	data := OneCurrData{Currency: currency}
	var date time.Time
	var value float64

	for rows.Next() {
		err := rows.Scan(&date, &value)
		if err != nil {
			return OneCurrData{}, err
		}

		data.Dates = append(data.Dates, date.Format("2006-01-02"))
		data.Values = append(data.Values, value)

	}

	if rows.Err() != nil {
		return OneCurrData{}, rows.Err()
	}

	return data, nil
}

type CurrenciesSymbols struct {
	Currencies []string `json:"currencies"`
}

func (d Database) SelectCurrenciesSymbolsV1(dateFrom string, dateTo string) (CurrenciesSymbols, error) {
	conn := d.connect()
	defer conn.Release()

	statement := `
    select
      distinct cs.symbol as curr_symbol
    from "data" d 
    left join "date" dt
    on d.date_id = dt.id 
    left join curr_symbol cs 
    on d.curr_symbol_id = cs.id 
    where dt.date between $1 and $2
    order by cs.symbol asc;
  `

	rows, err := conn.Query(context.Background(), statement, dateFrom, dateTo)
	if err != nil {
		return CurrenciesSymbols{}, err
	}
	defer rows.Close()

	data := CurrenciesSymbols{}
	currList := make([]string, 0, 130)
	symbol := ""

	for rows.Next() {
		err := rows.Scan(&symbol)
		if err != nil {
			return CurrenciesSymbols{}, err
		}
		currList = append(currList, symbol)
	}

	if rows.Err() != nil {
		return CurrenciesSymbols{}, rows.Err()
	}

	data.Currencies = currList

	return data, nil
}
