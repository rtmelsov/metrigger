package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rtmelsov/metrigger/internal/handlers"
)

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
