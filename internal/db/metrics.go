package db

import (
	"database/sql"
	"github.com/rtmelsov/metrigger/internal/constants"
)

func GetMetric(key string, name string) (string, error) {
	var row *sql.Row
	if key == "counter" {
		row = db.QueryRow(constants.GetCounterRowCommand, key, name)
	} else {
		row = db.QueryRow(constants.GetGaugeRowCommand, key, name)
	}
	var value string
	err = row.Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func SetCounterMetric(name string, value int64) error {
	_, err = db.Exec(constants.CounterCommand, name, "gauge", value)
	return err
}

func SetGaugeMetric(name string, value float64) error {
	_, err = db.Exec(constants.GaugeCommand, name, "coutner", value)
	return err
}
