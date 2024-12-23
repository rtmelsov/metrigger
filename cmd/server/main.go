package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[string]CounterMetric),
		gauge:   make(map[string]GaugeMetric),
	}
}

var mem = NewMemStorage()

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

func main() {
	err := run()
	if err != nil {
		log.Panic(err)
	}
}
func MetricsGaugeSet(name string, val string) error {
	met := NewGaugeMetric()
	met.Type = "gauge"
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	met.Value = f
	mem.SetGaugeMetric(name, *met)
	return nil
}

func MetricsCounterSet(name string, val string) error {
	met := NewCounterMetric()
	met.Type = "counter"
	i, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	oldCount := 0
	oldMet, err := mem.GetCounterMetric(name)
	if err != nil {
		fmt.Printf("Error: %v\r\n", err.Error())
	} else {
		oldCount = oldMet.Value
	}
	met.Value = oldCount + i
	mem.SetCounterMetric(name, *met)
	return nil
}

func run() error {
	fmt.Println("Server is running")
	return http.ListenAndServe(":8080", http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	paths := strings.Split(r.URL.String(), "/")
	if len(paths) == 5 {
		metType := paths[2]
		metName := paths[3]
		metVal := paths[4]
		switch metType {
		case "counter":
			err := MetricsCounterSet(metName, metVal)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case "gauge":
			err := MetricsGaugeSet(metName, metVal)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
