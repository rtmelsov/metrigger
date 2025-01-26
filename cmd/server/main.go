package main

import (
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	fmt.Println("server is running...")
	config.ServerParseFlag()

	logger := storage.GetMemStorage().GetLogger()

	prettyJSON, _ := json.MarshalIndent(config.ServerFlags, "", "  ")
	logger.Info("started", zap.String("server flags", string(prettyJSON)))

	defer logger.Sync()

	err := run()
	if err != nil {
		logger := storage.GetMemStorage().GetLogger()
		logger.Panic("error while running server", zap.String("error", err.Error()))
	}
}
func run() error {
	return http.ListenAndServe(config.ServerFlags.Addr, handlers.Webhook())
}
