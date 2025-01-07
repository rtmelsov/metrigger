package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/rtmelsov/metrigger/internal/handlers"
)

var data struct {
	Addr string `env:"ADDRESS"`
}

func ParseFlag() {

	flag.StringVar(&data.Addr, "a", "localhost:8080", "host and port to run server")

	flag.Parse()

	err := env.Parse(&data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ParseFlag()
	err := run()
	if err != nil {
		log.Panic(err)
	}
}
func run() error {
	fmt.Printf("Server is running: %v\r\n", data.Addr)
	return http.ListenAndServe(data.Addr, handlers.Webhook())
}
