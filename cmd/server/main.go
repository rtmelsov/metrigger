package server

import (
	"flag"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"log"
	"net/http"
)

var Addr string

func ParseFlag() {

	flag.StringVar(&Addr, "a", ":8080", "host and port to run server")

	flag.Parse()
}

func main() {
	ParseFlag()
	err := run()
	if err != nil {
		log.Panic(err)
	}
}
func run() error {
	fmt.Println("Server is running", Addr)
	return http.ListenAndServe(Addr, handlers.Webhook())
}
