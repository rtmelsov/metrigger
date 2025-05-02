package models

import "os"

// Metrics типа для работы с метриками - для записи/получение
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// ErrorType свой тип ошибки для отправки статуса через функции
type ErrorType struct {
	Text       string
	StatusCode int
}

// ServerFlagsType тип для получения данных флагов/env-переменных при запуске
type ServerFlagsType struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DataBaseDsn     string `env:"DATABASE_DSN"`
	JwtKey          string `env:"KEY"`
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
