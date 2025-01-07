package main

import (
	"fmt"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"log"
	"net/http"
)

func main() {
	err := run()
	if err != nil {
		log.Panic(err)
	}
}
func run() error {
	fmt.Println("Server is running")
	return http.ListenAndServe(":8080", http.HandlerFunc(handlers.Webhook))
}
