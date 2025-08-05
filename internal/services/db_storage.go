package services

import (
	"fmt"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/models"
	"strconv"
)

func GetDBCounter(name string) (*models.CounterMetric, error) {
	d, err := db.GetMetric("counter", name)
	if err != nil {
		return nil, err
	}
	var delta int64
	if d != "" {
		delta, err = strconv.ParseInt(d, 10, 64)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return nil, err
		}
	}

	return &models.CounterMetric{
		Type:  name,
		Value: delta,
	}, nil
}

func GetDBGauge(name string) (*models.GaugeMetric, error) {
	res, err := db.GetMetric("gauge", name)
	if err != nil {
		return nil, err
	}
	f, err := strconv.ParseFloat(res, 64)
	if err != nil {
		fmt.Println("Error converting string to float64:", err)
		return nil, err
	}

	return &models.GaugeMetric{
		Type:  name,
		Value: f,
	}, nil
}

func SetDBGauge(name string, val float64) error {

	return db.SetGaugeMetric(name, val)
}

func SetDBGounter(name string, val int64) error {
	return db.SetCounterMetric(name, val)
}
