package handlers

import (
	"net/http"

	"github.com/rtmelsov/metrigger/internal/constants"
	"github.com/rtmelsov/metrigger/internal/storage"
	"go.uber.org/zap"
	"html/template"
)

func MerticsListHandler(w http.ResponseWriter, r *http.Request) {
	mem := storage.GetMemStorage()
	t, err := template.New("webpage").Parse(constants.Tmpl)

	if err != nil {
		mem.GetLogger().Panic("Metric List Handler", zap.String("error", err.Error()))
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := t.Execute(w, mem); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
