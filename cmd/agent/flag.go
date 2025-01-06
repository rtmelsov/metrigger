package main

import (
	"flag"
)

var ReportInterval int
var PollInterval int

func ParseFlag() {

	flag.IntVar(&ReportInterval, "r", 10, "report interval")
	flag.IntVar(&PollInterval, "p", 2, "poll interval")

	flag.Parse()
}
