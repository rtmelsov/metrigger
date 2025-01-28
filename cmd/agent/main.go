package main

import (
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/agent"
	"github.com/rtmelsov/metrigger/internal/config"
	"go.uber.org/zap"
	"net"
	"time"
)

func waitForServer(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			return nil // Сервер доступен
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("services not available at %s after %v", address, timeout)
}

func main() {
	config.AgentParseFlag()
	logger := config.GetAgentStorage().GetLogger()

	// Проверка доступности сервера
	err := waitForServer(config.AgentFlags.Addr, 60*time.Second)
	if err != nil {
		logger.Error("Server not available", zap.String("error", err.Error()))
		return
	}

	prettyJSON, _ := json.MarshalIndent(config.AgentFlags, "", "  ")
	logger.Info("started", zap.String("agent flags", string(prettyJSON)))

	agent.Run()

}
