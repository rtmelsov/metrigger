package middleware

import (
	"bytes"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func JwtParser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storage.GetMemStorage().GetLogger().Info("jwt parser", zap.String("jwt", storage.ServerFlags.JwtKey))
		receivedHash := r.Header.Get("HashSHA256")
		if storage.ServerFlags.JwtKey != "" && receivedHash != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}

			expectedHash := helpers.ComputeHMACSHA256(body, storage.ServerFlags.JwtKey)

			storage.GetMemStorage().GetLogger().Info("comparing hashes",
				zap.String("expectedHash", expectedHash),
				zap.String("receivedHash", receivedHash))
			if receivedHash != expectedHash {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		h.ServeHTTP(w, r)
	})
}
