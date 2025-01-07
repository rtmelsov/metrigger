package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"

	"github.com/rtmelsov/metrigger/internal/handlers"
)

type Server struct {
	Addr string `env:"ADDRESS"`
}

func ServerFlags() *Server {
	var data Server

	flag.StringVar(&data.Addr, "a", ":8080", "host and port to run server")

	flag.Parse()
	err := env.Parse(&data)
	if err != nil {
		log.Fatal(err)
	}
	return &data
}

func main() {
	data := ServerFlags()
	err := run(*data)
	if err != nil {
		log.Panic(err)
	}
}
func run(data Server) error {
	fmt.Println("Server is running", data.Addr)
	return http.ListenAndServe(data.Addr, handlers.Webhook())
}
