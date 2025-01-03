package server

import (
	"fmt"
	"github.com/rtmelsov/metrigger/internal/storage"
	"strconv"
)

func MetricsGaugeSet(name string, val string) error {
	mem := storage.GetMem()
	met := storage.NewGaugeMetric()
	met.Type = "gauge"
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	met.Value = f
	mem.SetGaugeMetric(name, *met)
	return nil
}

func MetricsCounterGet(name string) (*storage.CounterMetric, error) {
	mem := storage.GetMem()
	oldMet, err := mem.GetCounterMetric(name)
	return oldMet, err
}

func MetricsGaugeGet(name string) (*storage.GaugeMetric, error) {
	mem := storage.GetMem()
	oldMet, err := mem.GetGaugeMetric(name)
	return oldMet, err
}

func MetricsGet() *storage.MemStorage {
	metrics := storage.GetMem()
	return metrics
}

func MetricsCounterSet(name string, val string) error {
	mem := storage.GetMem()
	met := storage.NewCounterMetric()
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
