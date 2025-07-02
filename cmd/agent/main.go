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

func waitForServer(address string) error {
	var timeouts = []int{1, 3, 5}
	for _, el := range timeouts {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			return nil // Сервер доступен
		}
		time.Sleep(time.Duration(el) * time.Second)
	}
	return fmt.Errorf("services not available at %s after %v", address, timeouts[len(timeouts)-1])
}

func main() {
	config.AgentParseFlag()
	logger := config.GetAgentStorage().GetLogger()

	// Проверка доступности сервера
	err := waitForServer(config.AgentFlags.Addr)
	if err != nil {
		logger.Error("Server not available", zap.String("error", err.Error()))
		return
	}

	prettyJSON, _ := json.MarshalIndent(config.AgentFlags, "", "  ")
	logger.Info("started", zap.String("agent flags", string(prettyJSON)))

	agent.Run()

}
