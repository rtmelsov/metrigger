package main

import (
	"context"
	"encoding/json"
	"fmt"
	// "github.com/rtmelsov/metrigger/cmd/staticlint"
	"github.com/rtmelsov/metrigger/internal/agent"
	"github.com/rtmelsov/metrigger/internal/config"
	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\r\n", buildVersion)
	fmt.Printf("Build date: %s\r\n", buildDate)
	fmt.Printf("Build commit: %s\r\n", buildCommit)
	// staticlint.Check()

	logger := config.GetAgentConfig().GetLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Проверка доступности сервера
	err := agent.WaitForServer(ctx, config.GetAgentConfig().Address())
	if err != nil {
		logger.Error("Server not available", zap.String("error", err.Error()))
		return
	}

	prettyJSON, err := json.MarshalIndent(config.GetAgentConfig().Address(), "", "  ")
	if err != nil {
		logger.Error("Error while try to marshal agent flags", zap.String("error", err.Error()))
		return
	}
	logger.Info("started", zap.String("agent flags", string(prettyJSON)))

	a := agent.NewAgent()

	a.Run()
}
