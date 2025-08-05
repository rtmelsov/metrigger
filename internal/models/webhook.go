package models

import "os"

// ErrorType свой тип ошибки для отправки статуса через функции
type ErrorType struct {
	Text       string
	StatusCode int
}

// ServerFlagsType тип для получения данных флагов/env-переменных при запуске
type ServerFlagsType struct {
	CryptoRate      string `env:"CRYPTO_KEY"`
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DataBaseDsn     string `env:"DATABASE_DSN"`
	JwtKey          string `env:"KEY"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
}

// CounterMetric тип для получения/записи значения counter в списке
type CounterMetric struct {
	Type  string
	Value int64
}

// GaugeMetric тип для получения/записи значения gauge в списке
type GaugeMetric struct {
	Type  string
	Value float64
}

// LocalStorage тип записи объекта в файл
type LocalStorage struct {
	Counter map[string]CounterMetric
	Gauge   map[string]GaugeMetric
	File    *os.File
}

// MetricsCollector тип записи списка в файл
type MetricsCollector map[string]float64

type MetricsCollectorData struct {
	Metrics *MetricsCollector
	Length  int
}

type ServerFileConfig struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DataBaseDsn   string `env:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
	TrustedSubnet string `json:"trusted_subnet"`
}
