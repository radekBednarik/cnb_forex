package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/radekBednarik/cnb_forex/data-getter/api"
	"github.com/radekBednarik/cnb_forex/data-getter/db"
	"github.com/radekBednarik/cnb_forex/data-getter/parser"
)

type Flags struct {
	configPath string
}

func flags() Flags {
	var f Flags

	f.configPath = *flag.String("c", "config.toml", "Provide file path to config.toml configuration file.")

	flag.Parse()

	return f
}

type Config struct {
	Date Date
}

type Date struct {
	Begin string
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file on path %s\n", path)
	}

	var config Config

	err = toml.Unmarshal(data, &config)

	if err != nil {
		log.Fatalf("Faild to unmarshal .toml data.\n%v\n", err)
	}

	return config, nil
}

func crunchData(dbs db.Database, dateBegin string) {
	fmt.Println("Data crunch started...")

	now := time.Now()
	dayDelta := 24 * time.Hour
	fNow := now.Format("02.01.2006")
	idCounter := 0

	for fNow != dateBegin {
		// call api to return data
		data, err := api.GetDailyData(fNow)
		if err != nil {
			log.Fatalf("Failed to get daily fx data for date '%s'. Error: %v", fNow, err)
		}

		// parse data
		parsedData := parser.ForexDataForDate{}
		parsedData.ParseFromText(data)

		fmt.Printf("Data point: %s\n", parsedData.Date)

		// check, if date is in db already, if yes, then continue
		// if we tried five consecutive previous dates and still its in the dbs
		// then break
		_, err = dbs.SelectIdFromTable(parsedData.Date, "date", "date")
		idCounter++
		if err == nil && idCounter >= 5 {
			break
		}

		// otherwise, insert data into db
		status := dbs.ProcessDailyData(&parsedData)
		if status {
			idCounter = 0
		}

		// adjust fNow to previous day
		now = now.Add(-dayDelta)
		fNow = now.Format("02.01.2006")

		// wait for a bit
		time.Sleep(200 * time.Millisecond)

		fmt.Println("=====")

	}

	fmt.Println("Done.")
}

func main() {
	flags := flags()
	config, err := loadConfig(flags.configPath)
	if err != nil {
		log.Fatalf("Failed to load config file.\n%v\n", err)
	}

	// connect to db and create connection pool
	connString := fmt.Sprintf("user=%s password=%s host=localhost port=5432 dbname=cnb_forex sslmode=verify-ca pool_max_conns=16", os.Getenv("USER"), os.Getenv("PASSWORD"))
	dbs := db.Database{}
	dbs.New(connString)
	defer dbs.Pool.Close()

	// check if tables are in db, if not, then create
	dbs.CreateTables()

	// process the data
	crunchData(dbs, config.Date.Begin)
}
