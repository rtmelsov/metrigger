package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rtmelsov/metrigger/internal/handlers"
)

func main() {
	err := run()
	if err != nil {
		log.Panic(err)
	}
}
func run() error {
	fmt.Println("Server is running")
	return http.ListenAndServe(":8080", handlers.Webhook())
}
