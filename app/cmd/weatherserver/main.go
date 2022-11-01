package main

import (
	"flag"
	"github.com/valen0k/wb_test_content/internal/app/weatherserver"
	"log"
)

var (
	//configPath       string
	locationFilePath string
)

func init() {
	//flag.StringVar(&configPath, "c", "configs/config.json", "path to config file")
	flag.StringVar(&locationFilePath, "l", "configs/locations.json", "path to locations file")
}

func main() {
	flag.Parse()

	log.Println("config initializing")
	config := weatherserver.GetConfig()

	if err := weatherserver.Start(config, locationFilePath); err != nil {
		log.Fatalln(err)
	}
}
