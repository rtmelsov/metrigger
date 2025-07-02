package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"sync"
)

type AgentFlagsType struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Addr           string `env:"ADDRESS"`
	JwtKey         string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

var (
	onceForAgent sync.Once
	AgentFlags   AgentFlagsType
	once         sync.Once
	agentMem     *AgentStorage
)

func AgentParseFlag() {
	onceForAgent.Do(func() {
		flag.IntVar(&AgentFlags.ReportInterval, "r", 10, "report interval")
		flag.StringVar(&AgentFlags.Addr, "a", "localhost:8080", "host and port to run services")
		flag.IntVar(&AgentFlags.PollInterval, "p", 2, "poll interval")
		flag.StringVar(&AgentFlags.JwtKey, "k", "server_key", "jwt key")
		flag.IntVar(&AgentFlags.RateLimit, "l", 5, "rate limit")

		flag.Parse()

		err := env.Parse(&AgentFlags)
		if err != nil {
			logger := storage.GetMemStorage().GetLogger()
			logger.Fatal("", zap.String("error", err.Error()))
		}
	})
}

type AgentStorage struct {
	Logger *zap.Logger
}

func (m *AgentStorage) GetLogger() *zap.Logger {
	return m.Logger
}

func GetAgentStorage() *AgentStorage {
	once.Do(func() {
		Log, _ := zap.NewProduction()

		agentMem = &AgentStorage{
			Logger: Log,
		}
	})
	return agentMem
}
