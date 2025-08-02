// Package config
package config

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/interfaces"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"log"
	"os"
	"sync"
)

var (
	once     sync.Once
	instance *models.AgentConfig
	flags    models.AgentFlags
)

// GetAgentConfig Создание и инициализация AgentConfig
func GetAgentConfig() interfaces.AgentActionsI {
	once.Do(func() {
		// CLI flags
		flags = models.AgentFlags{}
		flag.StringVar(&flags.ConfigCFile, "c", "", "")
		flag.StringVar(&flags.ConfigFile, "config", "", "")
		flag.StringVar(&flags.CryptoRate, "crypto-key", flags.CryptoRate, "")
		flag.IntVar(&flags.ReportInterval, "r", 10, "report interval")
		flag.StringVar(&flags.Addr, "a", "localhost:8080", "host and port")
		flag.IntVar(&flags.PollInterval, "p", 2, "poll interval")
		flag.StringVar(&flags.JwtKey, "k", "server_key", "jwt key")
		flag.IntVar(&flags.RateLimit, "l", 5, "rate limit")

		flag.Parse()

		if flags.ConfigFile != "" || flags.ConfigCFile != "" {
			var confs models.ClientFileConfig
			var conf string
			if flags.ConfigFile != "" {
				conf = flags.ConfigFile
			} else {
				conf = flags.ConfigCFile
			}
			data, err := os.ReadFile(conf)
			if err != nil {
				log.Fatal(err)
			}

			if err := json.Unmarshal(data, &confs); err != nil {
				log.Fatal(err)
			}

			if err = helpers.ClientFileConfigParser(&flags, &confs); err != nil {
				log.Fatal(err)
			}
		}

		// ENV
		if err := env.Parse(&flags); err != nil {
			logger := storage.GetMemStorage().GetLogger()
			logger.Fatal("failed to parse env", zap.Error(err))
		}

		// Logger
		logger, _ := zap.NewProduction()

		instance = &models.AgentConfig{
			Flags:  flags,
			Logger: logger,
		}
	})
	return instance
}
