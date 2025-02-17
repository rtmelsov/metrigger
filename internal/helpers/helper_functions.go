package helpers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/rtmelsov/metrigger/internal/models"
	"io"
	"net/http"
	"os"
)

// CompressData - Функция для упаковки данных в gzip
func CompressData(data []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}
	// Закрываем gzipWriter, чтобы завершить запись
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

// DecompressData - Функция для распаковки данных из gzip
func DecompressData(data []byte) (*bytes.Buffer, error) {
	buf := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	var result bytes.Buffer
	_, err = io.Copy(&result, gzipReader)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func JSONParse(r *http.Request) (*[]models.Metrics, error) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var metrics []models.Metrics // Slice to store users

	// Try parsing as a list (array)
	if err := json.Unmarshal(body, &metrics); err != nil {
		// If failed, try parsing as a single object
		var metric models.Metrics
		if err := json.Unmarshal(body, &metric); err != nil {
			return nil, err
		}
		// Convert single user to slice
		metrics = append(metrics, metric)
	}

	// If no users found
	if len(metrics) == 0 {
		return nil, err
	}

	return &metrics, nil
}

func EmptyLocalStorage(path string) (*models.LocalStorage, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &models.LocalStorage{
		Gauge:   make(map[string]models.GaugeMetric),
		Counter: make(map[string]models.CounterMetric),
		File:    file,
	}, nil
}
