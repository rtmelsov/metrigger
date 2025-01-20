package storage

import (
	"go.uber.org/zap"
	"sync"
)

var once sync.Once
var mem *MemStorage

type CounterMetric struct {
	Type  string
	Value int
}
type GaugeMetric struct {
	Type  string
	Value float64
}

type MemStorage struct {
	Counter map[string]CounterMetric
	Gauge   map[string]GaugeMetric
	mu      sync.Mutex
	logger  *zap.Logger
}

type MetricsStorage interface {
	GetGaugeMetric(name string) (*GaugeMetric, error)
	GetCounterMetric(name string) (*CounterMetric, error)
	SetGaugeMetric(name string, value GaugeMetric)
	SetCounterMetric(name string, value CounterMetric)
	GetLogger() *zap.Logger
}

func (m *MemStorage) GetLogger() *zap.Logger {
	return m.logger
}

func GetMemStorage() MetricsStorage {
	once.Do(func() {
		Log, _ := zap.NewProduction()
		mem = &MemStorage{
			Counter: make(map[string]CounterMetric),
			Gauge:   make(map[string]GaugeMetric),
			logger:  Log,
		}
	})
	return mem
}

func NewCounterMetric() *CounterMetric {
	return &CounterMetric{
		Type:  "",
		Value: 0,
	}
}

func NewGaugeMetric() *GaugeMetric {
	return &GaugeMetric{
		Type:  "",
		Value: 0,
	}
}

func (m *MemStorage) GetGaugeMetric(name string) (*GaugeMetric, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var value GaugeMetric
	value, ok := m.Gauge[name]
	if !ok {
		m.GetLogger().Info("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) GetCounterMetric(name string) (*CounterMetric, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var value CounterMetric
	value, ok := m.Counter[name]
	if !ok {
		m.GetLogger().Info("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) SetCounterMetric(name string, value CounterMetric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Counter[name] = value
}

func (m *MemStorage) SetGaugeMetric(name string, value GaugeMetric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Gauge[name] = value
}
