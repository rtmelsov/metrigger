package helpers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
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

func JsonParse(r *http.Request) (*models.Metrics, error) {
	if r.Body == nil {
		return nil, errors.New("body is empty")
	}
	var resp *models.Metrics
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&resp)

	if err != nil {
		return nil, err
	}
	if resp.ID == "" {
		return nil, errors.New("id is empty")
	}

	return resp, nil
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
