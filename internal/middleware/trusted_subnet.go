// Package middleware
package middleware

import (
	"github.com/rtmelsov/metrigger/internal/storage"
	"net/http"
)

func TrustedSubnet(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storage.GetMemStorage().GetLogger().Info("trusted subnet check")
		if r.Header.Get("X-Real-IP") != "" &&
			storage.ServerFlags.TrustedSubnet != "" {
			if r.Header.Get("X-Real-IP") != storage.ServerFlags.TrustedSubnet {
				http.Error(w, "x real ip is not correct", http.StatusForbidden)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
