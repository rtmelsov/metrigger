package main

import (
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	config.ServerParseFlag()

	logger := storage.GetMemStorage().GetLogger()
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
