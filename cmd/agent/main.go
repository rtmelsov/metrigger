package main

import (
	"github.com/rtmelsov/metrigger/internal/agent"
	"github.com/rtmelsov/metrigger/internal/config"
)

func main() {
	config.AgentParseFlag()
	agent.Run()
}
