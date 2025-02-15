package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
)

func ServerParseFlag() {
	storage.OnceForServer.Do(func() {
		flag.StringVar(&storage.ServerFlags.Addr, "a", "localhost:8080", "host and port to run services")
		flag.IntVar(&storage.ServerFlags.StoreInterval, "i", 300, "")
		flag.StringVar(&storage.ServerFlags.FileStoragePath, "f", "file.txt", "")
		flag.BoolVar(&storage.ServerFlags.Restore, "r", true, "")
		flag.StringVar(&storage.ServerFlags.DataBaseDsn, "d", "postgres://test:test@localhost:5432/dbname?sslmode=disable", "")

		flag.Parse()

		err := env.Parse(&storage.ServerFlags)
		if err != nil {
			logger := storage.GetMemStorage().GetLogger()
			logger.Fatal("error on parsing", zap.String("error", err.Error()))
		}
	})
}
