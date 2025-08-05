// Package agent
package agent

import (
	pb "github.com/rtmelsov/metrigger/proto"
)

// RequestToServer функция для формирования объекта для отправки в сервер
func RequestToServer(t string, key string, value float64, counter int64) *pb.Metric {
	return &pb.Metric{
		ID:    key,
		MType: t,
		Delta: counter,
		Value: value,
	}
}
