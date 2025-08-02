package config

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"os"
)

//postgres://postgres@localhost:5432/dbname?sslmode=disable

func ServerParseFlag() {
	storage.OnceForServer.Do(func() {
		flag.StringVar(&storage.ServerFlags.Addr, "a", "localhost:8080", "host and port to run services")
		flag.IntVar(&storage.ServerFlags.StoreInterval, "i", 300, "")
		flag.StringVar(&storage.ServerFlags.FileStoragePath, "f", "file.txt", "")
		flag.StringVar(&storage.ServerFlags.TrustedSubnet, "t", "", "")
		flag.StringVar(&storage.ServerFlags.CryptoRate, "crypto-key", "private.pem", "")
		flag.BoolVar(&storage.ServerFlags.Restore, "r", true, "")
		flag.StringVar(&storage.ServerFlags.DataBaseDsn, "d", "", "")
		flag.StringVar(&storage.ServerFlags.JwtKey, "k", "", "jwt key")

		flag.Parse()

		logger := storage.GetMemStorage().GetLogger()
		if flags.ConfigFile != "" || flags.ConfigCFile != "" {
			var confs models.ServerFileConfig
			var conf string
			if flags.ConfigFile != "" {
				conf = flags.ConfigFile
			} else {
				conf = flags.ConfigCFile
			}
			data, err := os.ReadFile(conf)
			if err != nil {
				logger.Fatal("error on parsing", zap.String("error", err.Error()))

			}

			if err := json.Unmarshal(data, &confs); err != nil {
				logger.Fatal("error on parsing", zap.String("error", err.Error()))

			}

			if err = helpers.ServerFileConfigParser(&storage.ServerFlags, &confs); err != nil {
				logger.Fatal("error on parsing", zap.String("error", err.Error()))
			}
		}

		err := env.Parse(&storage.ServerFlags)
		if err != nil {
			logger.Fatal("error on parsing", zap.String("error", err.Error()))
		}
	})
}
