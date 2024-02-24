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
	now := time.Now()
	dayDelta := 24 * time.Hour
	fNow := now.Format("01.01.2024")
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

		// check, if date is in db already, if yes, then continue
		// if we tried five consecutive previous dates and still its in the dbs
		// then break
		_, err = dbs.SelectIdFromTable(parsedData.Date, "date", "date")
		idCounter++
		if err == nil && idCounter >= 5 {
			break
		}
		idCounter = 0

		// otherwise, insert data into db
		dbs.ProcessDailyData(&parsedData)

		// adjust fNow to previous day
		now = now.Add(-dayDelta)
		fNow = now.Format("01.01.2024")

		// wait for a bit
		time.Sleep(200 * time.Millisecond)

	}
}

func main() {
	flags := flags()
	_, err := loadConfig(flags.configPath)
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

	// TODO: create a loop statement, which will go thru dates and stores
	// them to db. No concurrency, since we do not want to get banned by the server
	data, err := api.GetDailyData("17.01.2024")
	if err != nil {
		log.Fatalf("Attempt to GET daily cnb forex data failed with error:\n%v\n", err)
	}

	pData := parser.ForexDataForDate{}
	pData.ParseFromText(data)

	dbs.ProcessDailyData(&pData)

	fmt.Println(pData)
}
