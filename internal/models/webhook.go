package models

import "os"

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type ErrorType struct {
	Text       string
	StatusCode int
}

type ServerFlagsType struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DataBaseDsn     string `env:"DATABASE_DSN"`
}

type CounterMetric struct {
	Type  string
	Value int
}
type GaugeMetric struct {
	Type  string
	Value float64
}

type LocalStorage struct {
	Counter map[string]CounterMetric
	Gauge   map[string]GaugeMetric
	File    *os.File
}

type MetricsCollector map[string]float64
