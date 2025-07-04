package main

import (
	"context"
	"encoding/json"
	"github.com/rtmelsov/metrigger/cmd/staticlint"
	"github.com/rtmelsov/metrigger/internal/agent"
	"github.com/rtmelsov/metrigger/internal/config"
	"go.uber.org/zap"
)

func main() {
	staticlint.Check()

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

	agent.Run()
}
