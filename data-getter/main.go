package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

func main() {
	flags := flags()
	_, err := loadConfig(flags.configPath)
	if err != nil {
		log.Fatalf("Failed to load config file.\n%v\n", err)
	}

	// TODO: here - place db tables creation if no tables found in db
	connString := fmt.Sprintf("user=%s password=%s host=localhost port=5432 dbname=cnb_forex sslmode=verify-ca pool_max_conns=16", os.Getenv("USER"), os.Getenv("PASSWORD"))

	dbs := db.Database{}
	dbs.New(connString)
	db.CreateTables(dbs)

	// TODO: next - start retrieving data from now() to config.date.begin
	// each date step checks, if it is in the database. if not, continue, if yes, break
	// since data are here and so we do not need to continue
	// first load however will be db calls heavy

	data, err := api.GetDailyData("16.01.2024")
	if err != nil {
		log.Fatalf("Attempt to GET daily cnb forex data failed with error:\n%v\n", err)
	}

	pData := parser.ForexDataForDate{}
	pData.ParseFromText(data)

	fmt.Println(pData)
}
