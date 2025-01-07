package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rtmelsov/metrigger/internal/handlers"
	"github.com/rtmelsov/metrigger/internal/handlers/helpers"
	"github.com/rtmelsov/metrigger/internal/models"
)

func main() {
	data := helpers.ServerFlags()
	err := run(*data)
	if err != nil {
		log.Panic(err)
	}
}
func run(data models.ServerFlags) error {
	fmt.Println("Server is running", data.Addr)
	return http.ListenAndServe(data.Addr, handlers.Webhook())
}
