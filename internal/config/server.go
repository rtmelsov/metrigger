package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"sync"
)

type ServerFlagsType struct {
	Addr string `env:"ADDRESS"`
}

var (
	onceForServer sync.Once
	ServerFlags   ServerFlagsType
)

func ServerParseFlag() {
	onceForServer.Do(func() {
		flag.StringVar(&ServerFlags.Addr, "a", "localhost:8080", "host and port to run services")

		flag.Parse()

		err := env.Parse(&ServerFlags)
		if err != nil {
			logger := storage.GetMemStorage().GetLogger()
			logger.Fatal("error on parsing", zap.String("error", err.Error()))
		}
	})
}
