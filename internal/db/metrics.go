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

func SetMetric(key string, name string, value string) error {
	if key == "counter" {
		_, err = db.Exec(constants.CounterCommand, name, key, value)
	} else {
		_, err = db.Exec(constants.GaugeCommand, name, key, value)
	}
	return err
}
