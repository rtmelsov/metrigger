package handlers

import (
	"database/sql"
	"errors"
	"github.com/rtmelsov/metrigger/internal/constants"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func PingDBHandler(w http.ResponseWriter, r *http.Request) {
	mem := storage.GetMemStorage()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	database, err := db.GetDataBase()
	if err != nil {
		mem.GetLogger().Panic("Error to ping to db", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = database.Ping()
	if err != nil {
		mem.GetLogger().Panic("Error to ping to db", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("ok")); err != nil {
		mem.GetLogger().Panic("Error while sending response", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateMetrics(response *[]models.Metrics) (*[]models.Metrics, error) {
	Log := storage.GetMemStorage().GetLogger()

	Log.Info("UpdateMetrics 1")
	var newMetrics []models.Metrics
	DB, err := db.GetDataBase()
	if err != nil {
		return nil, err
	}

	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	setGauge, setCounter, getGommand, err := getCommands(tx)
	if err != nil {
		Log.Panic("error while get command", zap.Error(err))
		return nil, err
	}

	Log.Info("UpdateMetrics 2")
	defer setGauge.Close()
	defer setCounter.Close()
	defer getGommand.Close()

	for _, v := range *response {

		switch v.MType {
		case "gauge":
			if v.Value == nil {
				return nil, errors.New("value is empty")
			}
			_, err = setGauge.Exec(v.ID, v.MType, v.Value)

		case "counter":
			if v.Delta == nil {
				return nil, errors.New("delta is empty")
			}
			_, err = setCounter.Exec(v.ID, v.MType, v.Delta)

		default:

		}

		if err != nil {
			tx.Rollback()
			Log.Panic("error while exec", zap.Error(err))
			return nil, err
		}
		var value float64
		var delta int64
		if v.MType == "gauge" {
			err = getGommand.QueryRow(v.MType, v.ID).Scan(&value)
			newMetrics = append(newMetrics, models.Metrics{
				ID:    v.ID,
				MType: v.MType,
				Value: &value,
			})
		} else {
			err = getGommand.QueryRow(v.MType, v.ID).Scan(&delta)
			newMetrics = append(newMetrics, models.Metrics{
				ID:    v.ID,
				MType: v.MType,
				Delta: &delta,
			})
		}
	}

	Log.Info("UpdateMetrics 2")
	err = tx.Commit()
	if err != nil {
		Log.Panic("error while commit", zap.Error(err))
		return nil, err
	}
	return &newMetrics, nil
}

func getCommands(tx *sql.Tx) (*sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	setGauge, err := tx.Prepare(constants.GaugeCommand)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	setCounter, err := tx.Prepare(constants.CounterCommand)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	getCommand, err := tx.Prepare(constants.GetRowCommand)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	return setGauge, setCounter, getCommand, nil
}
