package agent

import (
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/metrigger/internal/config"
	"github.com/rtmelsov/metrigger/internal/helpers"
	"github.com/rtmelsov/metrigger/internal/models"
	"go.uber.org/zap"
	"net/http"
)

// Worker функция для отправки POST запроса
func Worker(metricList []*models.Metrics) error {
	data, err := json.Marshal(&metricList)
	logger := config.GetAgentConfig().GetLogger()
	if err != nil {
		logger.Error("Error to Marshal JSON", zap.String("error", err.Error()))
		return err
	}

	reqBody, err := helpers.CompressData(data)
	if err != nil {
		logger.Error("Error to Marshal JSON", zap.String("error", err.Error()))
		return err
	}

	url := fmt.Sprintf("http://%s/updates/", config.GetAgentConfig().Address())

	req, err := http.NewRequest("POST", url, reqBody)

	if err != nil {
		logger.Error("1 Request to services", zap.String("error", err.Error()))
		return err
	}

	if config.GetAgentConfig().JwtKey() != "" {
		hash := helpers.ComputeHMACSHA256(data, config.GetAgentConfig().JwtKey())
		req.Header.Set("HashSHA256", hash)
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logger.Error("2 Request to services", zap.String("error", err.Error()))
		return err
	}

	if resp.StatusCode != 200 {
		logger.Error("status code: ", zap.Int("code", resp.StatusCode))
	}

	if config.GetAgentConfig().JwtKey() != "" {
		logger.Info("hash answer", zap.String("HashSHA256", resp.Header.Get("HashSHA256")))
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Error("3 Request to services", zap.String("error", err.Error()))
		return err
	}
	logger.Info("requested")
	return nil
}
