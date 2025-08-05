package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"os"
	"sync"
)

type MemStorage struct {
	Counter map[string]models.CounterMetric
	Gauge   map[string]models.GaugeMetric
	File    *os.File
	Mu      sync.Mutex
	Logger  *zap.Logger
	Writer  *bufio.Writer
}

var (
	once          sync.Once
	serverMem     *MemStorage
	OnceForServer sync.Once
	ServerFlags   models.ServerFlagsType
)

func (m *MemStorage) GetLogger() *zap.Logger {
	return m.Logger
}

func GetMemStorage() *MemStorage {
	once.Do(func() {
		Log, _ := zap.NewProduction()
		file, err := getDataFromFile()
		if err != nil {
			Log.Error(err.Error())
		}
		serverMem = &MemStorage{
			Counter: file.Counter,
			Gauge:   file.Gauge,
			File:    file.File,
			Logger:  Log,
		}
	})
	return serverMem
}

func NewCounterMetric() *models.CounterMetric {
	return &models.CounterMetric{
		Type:  "",
		Value: 0,
	}
}

func NewGaugeMetric() *models.GaugeMetric {
	return &models.GaugeMetric{
		Type:  "",
		Value: 0,
	}
}

func (m *MemStorage) GetGaugeMetric(name string) (*models.GaugeMetric, error) {
	value, ok := m.Gauge[name]
	if !ok {
		return nil, errors.New("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) GetCounterMetric(name string) (*models.CounterMetric, error) {
	value, ok := m.Counter[name]
	if !ok {
		return nil, errors.New("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) SetCounterMetric(name string, value models.CounterMetric) {
	m.Counter[name] = value
	if err := m.WriteAllData(ServerFlags.StoreInterval); err != nil {
		m.Logger.Error(err.Error())
	}
}

func (m *MemStorage) SetGaugeMetric(name string, value models.GaugeMetric) {
	m.Gauge[name] = value
	if err := m.WriteAllData(ServerFlags.StoreInterval); err != nil {
		m.Logger.Error(err.Error())
	}
}

func (m *MemStorage) SetDataToFile(value models.GaugeMetric) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := m.Writer.Write(data); err != nil {
		return err
	}

	// записываем буфер в файл
	return m.Writer.Flush()
}
