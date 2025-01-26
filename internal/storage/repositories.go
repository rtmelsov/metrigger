package storage

import (
	"encoding/json"
	"errors"
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

		prettyJSON, _ := json.MarshalIndent(mem, "", "  ")
		Log.Info("first time:",
			zap.String("mem - ", string(prettyJSON)))
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
	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("get data:",
		zap.String("GetGaugeMetric name", name),
		zap.String("GetGaugeMetric value", string(prettyJSON)))
	if !ok {
		return nil, errors.New("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) GetCounterMetric(name string) (*CounterMetric, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var value CounterMetric
	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("get data:",
		zap.String("GetCounterMetric name", name),
		zap.String("GetCounterMetric value", string(prettyJSON)))
	value, ok := m.Counter[name]
	if !ok {
		return nil, errors.New("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) SetCounterMetric(name string, value CounterMetric) {
	m.mu.Lock()
	defer m.mu.Unlock()

	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("set data:",
		zap.String("SetCounterMetric name", name),
		zap.String("SetCounterMetric value", string(prettyJSON)))

	m.Counter[name] = value
}

func (m *MemStorage) SetGaugeMetric(name string, value GaugeMetric) {
	m.mu.Lock()
	defer m.mu.Unlock()

	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("set data:",
		zap.String("set gauge metric name", name),
		zap.String("set gauge metric value", string(prettyJSON)))

	m.Gauge[name] = value
}
