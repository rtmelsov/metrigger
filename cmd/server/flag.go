package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

var Data struct {
	Addr string `env:"ADDRESS"`
}

func ParseFlag() {

	flag.StringVar(&Data.Addr, "a", ":8080", "host and port to run server")

	flag.Parse()
	err := env.Parse(&Data)
	if err != nil {
		log.Fatal(err)
	}
}
