package models

import "go.uber.org/zap"

type AgentFlags struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Addr           string `env:"ADDRESS"`
	JwtKey         string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

type AgentConfig struct {
	Flags  AgentFlags
	Logger *zap.Logger
}

func (cfg *AgentConfig) ReportInterval() int {
	return cfg.Flags.ReportInterval
}

func (cfg *AgentConfig) PollInterval() int {
	return cfg.Flags.PollInterval
}

func (cfg *AgentConfig) Address() string {
	return cfg.Flags.Addr
}

func (cfg *AgentConfig) JwtKey() string {
	return cfg.Flags.JwtKey
}

func (cfg *AgentConfig) RateLimit() int {
	return cfg.Flags.RateLimit
}

func (cfg *AgentConfig) GetLogger() *zap.Logger {
	return cfg.Logger
}
