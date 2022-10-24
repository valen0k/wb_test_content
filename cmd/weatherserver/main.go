package main

import (
	"flag"
	"github.com/valen0k/wb_test_content/internal/app/weatherserver"
	"log"
)

var (
	configPath       string
	locationFilePath string
)

func init() {
	flag.StringVar(&configPath, "c", "configs/config.json", "path to config file")
	flag.StringVar(&locationFilePath, "l", "configs/locations.json", "path to locations file")
}

func main() {
	flag.Parse()

	config := weatherserver.NewConfig()
	if err := config.DecodeConfigFile(configPath); err != nil {
		log.Fatalln(err)
	}

	server := weatherserver.New(config)
	if err := server.Start(locationFilePath); err != nil {
		log.Fatalln(err)
	}
}
