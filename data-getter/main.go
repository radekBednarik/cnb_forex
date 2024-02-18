package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml"
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
	tomlConfig, err := loadConfig(flags.configPath)

	if err != nil {
		log.Fatalf("Failed to load config file.\n%v\n", err)
	}

	fmt.Println(tomlConfig.Date.Begin)

}
