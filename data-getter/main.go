package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
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

func loadConfig(path string) (toml.MetaData, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("Failed to read file on path %s\n", path)
	}

	type date struct {
		begin string
	}

	type config map[string]date

	var con config

	pConf, err := toml.Decode(string(data), &con)

	if err != nil {
		log.Fatalf("Failed to decode %s to TOML datastructure.", string(data))
	}

	return pConf, nil
}

func main() {
	flags := flags()
	tomlConfig, err := loadConfig(flags.configPath)

	if err != nil {
		log.Fatalf("Failed to load config file.\n%v\n", err)
	}

	fmt.Println(tomlConfig.Undecoded())
}
