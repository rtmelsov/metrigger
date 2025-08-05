// Package storage
package storage

import (
	"encoding/json"
	"errors"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/models"
	"os"
)

func SetDataToFile(m *MemStorage) error {
	if ServerFlags.FileStoragePath == "" {
		return errors.New("file storage path is not exist")
	}

	data := map[string]any{
		"GaugeMetric":   m.Gauge,
		"CounterMetric": m.Counter,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Открываем файл в режиме записи с очисткой
	file, err := os.OpenFile(ServerFlags.FileStoragePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	return err
}

func (m *MemStorage) WriteAllData(n int) error {
	if n == 0 {
		return SetDataToFile(m)
	}
	go func() {
		m.Mu.Lock()
		defer m.Mu.Unlock()
		_ = SetDataToFile(m)
	}()
	return nil
}

func getDataFromFile() (*models.LocalStorage, error) {
	if !ServerFlags.Restore {
		return helpers.EmptyLocalStorage(ServerFlags.FileStoragePath)
	}
	file, err := os.Open(ServerFlags.FileStoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			return helpers.EmptyLocalStorage(ServerFlags.FileStoragePath)
		}
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() < 1 {
		return helpers.EmptyLocalStorage(ServerFlags.FileStoragePath)
	}

	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	var result map[string]map[string]any

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	gauge := make(map[string]models.GaugeMetric)
	counter := make(map[string]models.CounterMetric)

	if rawData1, ok := result["GaugeMetric"]; ok {
		for key, value := range rawData1 {
			valueBytes, _ := json.Marshal(value)
			var data1 models.GaugeMetric
			err = json.Unmarshal(valueBytes, &data1)
			if err != nil {
				return nil, err
			}
			gauge[key] = data1
		}
	}

	if rawData2, ok := result["CounterMetric"]; ok {
		for key, value := range rawData2 {
			valueBytes, _ := json.Marshal(value)
			var data2 models.CounterMetric
			err = json.Unmarshal(valueBytes, &data2)
			if err != nil {
				return nil, err
			}
			counter[key] = data2
		}
	}
	var localstorage = models.LocalStorage{
		Gauge:   gauge,
		Counter: counter,
		File:    file,
	}

	return &localstorage, nil
}
