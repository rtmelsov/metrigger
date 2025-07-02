package handlers

import (
	"github.com/rtmelsov/metrigger/internal/db"
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
