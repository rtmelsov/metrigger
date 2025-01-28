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
}

var (
	onceForAgent sync.Once
	AgentFlags   AgentFlagsType
)

func AgentParseFlag() {
	onceForAgent.Do(func() {
		flag.IntVar(&AgentFlags.ReportInterval, "r", 10, "report interval")
		flag.StringVar(&AgentFlags.Addr, "a", "localhost:8080", "host and port to run services")
		flag.IntVar(&AgentFlags.PollInterval, "p", 2, "poll interval")

		flag.Parse()

		err := env.Parse(&AgentFlags)
		if err != nil {
			logger := storage.GetMemStorage().GetLogger()
			logger.Fatal("", zap.String("error", err.Error()))
		}
	})
}
