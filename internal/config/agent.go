package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rtmelsov/metrigger/internal/interfaces"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"sync"
)

var (
	once     sync.Once
	instance *models.AgentConfig
)

// GetAgentConfig Создание и инициализация AgentConfig
func GetAgentConfig() interfaces.AgentActionsI {
	once.Do(func() {
		flags := models.AgentFlags{
			ReportInterval: 10,
			PollInterval:   2,
			Addr:           "localhost:8080",
			JwtKey:         "server_key",
			RateLimit:      5,
		}

		// CLI flags
		flag.IntVar(&flags.ReportInterval, "r", flags.ReportInterval, "report interval")
		flag.StringVar(&flags.Addr, "a", flags.Addr, "host and port")
		flag.IntVar(&flags.PollInterval, "p", flags.PollInterval, "poll interval")
		flag.StringVar(&flags.JwtKey, "k", flags.JwtKey, "jwt key")
		flag.IntVar(&flags.RateLimit, "l", flags.RateLimit, "rate limit")
		flag.Parse()

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
