package handlers

import (
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"net/http"
)

func SendData(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	// Если ключ есть, добавляем хеш в заголовок ответа
	if storage.ServerFlags.JwtKey != "" {
		responseHash := helpers.ComputeHMACSHA256(data, storage.ServerFlags.JwtKey)
		w.Header().Set("HashSHA256", responseHash)
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}
