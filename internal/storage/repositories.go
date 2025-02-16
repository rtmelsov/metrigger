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

type ReadMetricsStorage interface {
	GetGaugeMetric(name string) (*models.GaugeMetric, error)
	GetCounterMetric(name string) (*models.CounterMetric, error)
	GetLogger() *zap.Logger
}

type SetMetricStorage interface {
	SetGaugeMetric(name string, value models.GaugeMetric)
	SetCounterMetric(name string, value models.CounterMetric)
	SetDataToFile(value models.CounterMetric) error
}

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

		prettyJSON, _ := json.MarshalIndent(serverMem, "", "  ")
		Log.Info("first time:",
			zap.String("mem - ", string(prettyJSON)))
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
	var value models.GaugeMetric
	value, ok := m.Gauge[name]
	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("get data:",
		zap.String("GetGaugeMetric name", name),
		zap.String("GetGaugeMetric value", string(prettyJSON)))
	if !ok {
		return nil, errors.New("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) GetCounterMetric(name string) (*models.CounterMetric, error) {
	var value models.CounterMetric
	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("get data:",
		zap.String("GetCounterMetric name", name),
		zap.String("GetCounterMetric value", string(prettyJSON)))
	value, ok := m.Counter[name]
	if !ok {
		return nil, errors.New("can't get that name's value")
	}
	return &value, nil
}

func (m *MemStorage) SetCounterMetric(name string, value models.CounterMetric) {
	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("set data:",
		zap.String("SetCounterMetric name", name),
		zap.String("SetCounterMetric value", string(prettyJSON)))

	m.Counter[name] = value
	if err := m.WriteAllData(ServerFlags.StoreInterval); err != nil {
		m.Logger.Error(err.Error())
	}
}

func (m *MemStorage) SetGaugeMetric(name string, value models.GaugeMetric) {
	logger := GetMemStorage().GetLogger()
	prettyJSON, _ := json.MarshalIndent(value, "", "  ")
	logger.Info("set data:",
		zap.String("set gauge metric name", name),
		zap.String("set gauge metric value", string(prettyJSON)))
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
