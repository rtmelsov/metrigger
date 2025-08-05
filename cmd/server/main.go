package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/handlers"
	"github.com/rtmelsov/metrigger/internal/middleware"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
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

	run(logger)
}
func run(logger *zap.Logger) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		logger.Info("error to try to listen 3200 port", zap.String("error", err.Error()))
		return
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Logger,
			middleware.TrustedSubnet,
		),
	)
	handlers.InitWebhook(s)

	idleConnsClosed := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-quit
		s.Stop()
		close(idleConnsClosed)
	}()

	if err := s.Serve(listen); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("error while running services", zap.String("error", err.Error()))
	}
	<-idleConnsClosed
	logger.Info("server exiting")
}
