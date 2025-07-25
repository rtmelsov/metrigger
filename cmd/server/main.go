package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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

	// staticlint.Check()
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

	defer logger.Sync()

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

	run(logger)
}
func run(logger *zap.Logger) {
	srv := &http.Server{
		Addr:    storage.ServerFlags.Addr,
		Handler: handlers.Webhook(),
	}
	srv.SetKeepAlivesEnabled(false)

	idleConnsClosed := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("error while shutdown", zap.String("error", err.Error()))
		} else {
			logger.Info("shutdown complete")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("error while running services", zap.String("error", err.Error()))
	}
	<-idleConnsClosed
	logger.Info("server exiting")
}
