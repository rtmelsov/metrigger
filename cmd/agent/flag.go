package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
)

var Env struct {
	ReportInterval int      `env:"REPORT_INTERVAL"`
	PollInterval   int      `env:"POLL_INTERVAL"`
	Address        []string `env:"ADDRESS" envSeparator:":"`
}

func FlagParse() {
	flag.IntVar(&Env.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&Env.PollInterval, "p", 2, "poll interval")
	Env.Address = []string{"localhost", "8080"}

	flag.Parse()

	err := env.Parse(&Env)
	fmt.Println(Env)
	if err != nil {
		log.Fatal(err)
	}
}
