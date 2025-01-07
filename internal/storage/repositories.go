package storage

import (
	"errors"
	"sync"
)

type CounterMetric struct {
	Type  string
	Value int
}
type GaugeMetric struct {
	Type  string
	Value float64
}

type MemStorage struct {
	counter map[string]CounterMetric
	gauge   map[string]GaugeMetric
	mu      sync.Mutex
}

func newMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[string]CounterMetric),
		gauge:   make(map[string]GaugeMetric),
	}
}

var mem = newMemStorage()

func GetMem() *MemStorage {
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

func (m *MemStorage) GetCounterMetric(name string) (*CounterMetric, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var value CounterMetric
	value, ok := m.counter[name]
	if !ok {
		return nil, errors.New("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) SetCounterMetric(name string, value CounterMetric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[name] = value
}

func (m *MemStorage) SetGaugeMetric(name string, value GaugeMetric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[name] = value
}
