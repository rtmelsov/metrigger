package handlers

import (
	"database/sql"
	"github.com/rtmelsov/metrigger/internal/constants"
	"github.com/rtmelsov/metrigger/internal/db"
	"github.com/rtmelsov/metrigger/internal/models"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

// PingDBHandler обрабатывает HTTP-запрос для проверки соединения с базой данных.
//
// При успешной проверке возвращает статус 200 OK и сообщение "ok".
// В случае ошибки — пишет ошибку в ответ и логирует её.
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

// UpdateMetrics обновляет метрики в базе данных.
//
// Принимает список метрик `response`, обновляет их значения в транзакции и возвращает
// актуальные значения метрик. В случае ошибки возвращает ошибку.
//
// Возможные ошибки:
// - ошибка подключения к базе данных
// - ошибка выполнения SQL-запросов
// - ошибка при парсинге данных
func UpdateMetrics(response *[]models.Metrics) (*[]models.Metrics, error) {
	log := storage.GetMemStorage().GetLogger()

	tx, err := db.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer db.RollbackTx(tx)

	setGauge, setCounter, getGaugeCommand, getDeltaCommand, err := getCommands(tx)
	if err != nil {
		log.Panic("get command error", zap.Error(err))
		return nil, err
	}
	defer db.CloseStmt(setGauge)
	defer db.CloseStmt(setCounter)
	defer db.CloseStmt(getGaugeCommand)
	defer db.CloseStmt(getDeltaCommand)

	var newMetrics []models.Metrics
	for _, v := range *response {
		log.Info("processing metric", zap.String("type", v.MType), zap.String("id", v.ID))
		if err := db.InsertMetric(v, setGauge, setCounter); err != nil {
			log.Panic("insert metric error", zap.Error(err))
			return nil, err
		}

		updated, err := db.FetchUpdatedMetric(v, getGaugeCommand, getDeltaCommand)
		if err != nil {
			log.Panic("fetch updated metric error", zap.Error(err))
			return nil, err
		}
		newMetrics = append(newMetrics, updated)
	}

	if err := tx.Commit(); err != nil {
		log.Panic("commit error", zap.Error(err))
		return nil, err
	}
	return &newMetrics, nil
}

func getCommands(tx *sql.Tx) (*sql.Stmt, *sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	setGauge, err := tx.Prepare(constants.GaugeCommand)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	setCounter, err := tx.Prepare(constants.CounterCommand)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	getGaugeCommand, err := tx.Prepare(constants.GetGaugeRowCommand)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	getCounterCommand, err := tx.Prepare(constants.GetCounterRowCommand)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return setGauge, setCounter, getGaugeCommand, getCounterCommand, nil
}
