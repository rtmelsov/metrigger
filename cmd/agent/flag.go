package main

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"log"
)

var Env struct {
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

func FlagParse() {
	flag.IntVar(&Env.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&Env.PollInterval, "p", 2, "poll interval")

	flag.Parse()

	err := env.Parse(&Env)
	if err != nil {
		log.Fatal(err)
	}
}
