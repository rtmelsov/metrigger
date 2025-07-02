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

func worker(metricList []*models.Metrics) {
	data, err := json.Marshal(&metricList)
	logger := config.GetAgentStorage().GetLogger()
	if err != nil {
		logger.Panic("Error to Marshal JSON", zap.String("error", err.Error()))
		return
	}

	reqBody, err := helpers.CompressData(data)
	if err != nil {
		logger.Panic("Error to Marshal JSON", zap.String("error", err.Error()))
		return
	}

	url := fmt.Sprintf("http://%s/updates/", config.AgentFlags.Addr)

	req, err := http.NewRequest("POST", url, reqBody)

	if err != nil {
		logger.Panic("1 Request to services", zap.String("error", err.Error()))
		return
	}

	if config.AgentFlags.JwtKey != "" {
		hash := helpers.ComputeHMACSHA256(data, config.AgentFlags.JwtKey)
		req.Header.Set("HashSHA256", hash)
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logger.Panic("2 Request to services", zap.String("error", err.Error()))
		return
	}

	if resp.StatusCode != 200 {
		logger.Error("status code: ", zap.Int("code", resp.StatusCode))
	}

	if config.AgentFlags.JwtKey != "" {
		logger.Info("hash answer", zap.String("HashSHA256", resp.Header.Get("HashSHA256")))
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Panic("3 Request to services", zap.String("error", err.Error()))
		return
	}
	logger.Info("requested")
}
