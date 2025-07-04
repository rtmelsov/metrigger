package main

import (
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/cmd/staticlint"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/pprof"
	"time"
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

	staticlint.Check()
	config.ServerParseFlag()

	logger := storage.GetMemStorage().GetLogger()

	if storage.ServerFlags.DataBaseDsn != "" {
		_, err := db.GetDataBase()
		if err != nil {
			logger.Panic("error while running services", zap.String("error", err.Error()))
			return
		} else {
			logger.Info("database is connected", zap.String("timestamp", time.Now().Format(time.RFC3339)))
		}
	}

	prettyJSON, _ := json.MarshalIndent(storage.ServerFlags, "", "  ")
	logger.Info("started", zap.String("services flags", string(prettyJSON)))

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error(err.Error())
		}
	}(logger)

	go func(logger *zap.Logger) {
		pprofMux := http.NewServeMux()
		pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		logger.Info("start pprof on 6060")

		err := http.ListenAndServe("localhost:6060", pprofMux)
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	err := run()
	if err != nil {
		logger.Panic("error while running services", zap.String("error", err.Error()))
	}

}
func run() error {
	return http.ListenAndServe(storage.ServerFlags.Addr, handlers.Webhook())
}
