package services

import (
	"fmt"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/models"
	"strconv"
)

func GetDBCounter(name string) (*models.CounterMetric, error) {
	res, err := db.GetMetric("counter", name)
	if err != nil {
		return nil, err
	}
	num, err := strconv.Atoi(res)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return nil, err
	}

	return &models.CounterMetric{
		Type:  name,
		Value: num,
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

func SetDBGauge(name string, val string) error {
	return db.SetMetric("gauge", name, val)
}

func SetDBGounter(name string, val string) error {
	return db.SetMetric("counter", name, val)
}
