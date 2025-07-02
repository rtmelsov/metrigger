package main

import (
	"encoding/json"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	config.ServerParseFlag()

	logger := storage.GetMemStorage().GetLogger()

	if storage.ServerFlags.DataBaseDsn != "" {
		_, err := db.GetDataBase()
		if err != nil {
			logger.Panic("error while running services", zap.String("error", err.Error()))
			return
		}
	}

	prettyJSON, _ := json.MarshalIndent(storage.ServerFlags, "", "  ")
	logger.Info("started", zap.String("services flags", string(prettyJSON)))

	defer logger.Sync()

	err := run()
	if err != nil {
		logger.Panic("error while running services", zap.String("error", err.Error()))
		return
	}
}
func run() error {
	return http.ListenAndServe(storage.ServerFlags.Addr, handlers.Webhook())
}
