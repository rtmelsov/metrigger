package main

import (
	"flag"
)

var Addr string

func ParseFlag() {

	flag.StringVar(&Addr, "a", ":8080", "host and port to run server")

	flag.Parse()
}
