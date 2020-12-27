package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/jpurdie/authapi/pkg/api"
	"github.com/jpurdie/authapi/pkg/utl/config"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
