package main

import (
	"github.com/rtmelsov/metrigger/internal/agent"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"net/http"
)

func main() {
	config.AgentParseFlag()
	go func() error {
		return http.ListenAndServe(config.AgentFlags.Addr, handlers.Webhook())
	}()
	agent.Run()
}
