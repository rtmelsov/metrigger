// Package services
package services

import (
	"github.com/rtmelsov/metrigger/internal/storage"
	pb "github.com/rtmelsov/metrigger/proto"
)

func MetricsGaugeSet(resp *pb.Metric) error {
	mem := storage.GetMemStorage()
	mem.Mu.Lock()
	defer mem.Mu.Unlock()
	if storage.ServerFlags.DataBaseDsn != "" {
		return SetDBGauge(resp.ID, resp.Value)
	}
	met := storage.NewGaugeMetric()
	met.Type = "gauge"
	met.Value = resp.Value
	mem.SetGaugeMetric(resp.ID, *met)
	return nil
}

func MetricsCounterSet(resp *pb.Metric) error {
	mem := storage.GetMemStorage()
	mem.Mu.Lock()
	defer mem.Mu.Unlock()
	if storage.ServerFlags.DataBaseDsn != "" {
		return SetDBGounter(resp.ID, resp.Delta)
	}
	met := storage.NewCounterMetric()
	met.Type = "counter"
	var oldCount int64 = 0
	oldMet, err := mem.GetCounterMetric(resp.ID)
	if err == nil {
		oldCount = oldMet.Value
	}
	met.Value = oldCount + resp.Delta
	mem.SetCounterMetric(resp.ID, *met)
	return nil
}
