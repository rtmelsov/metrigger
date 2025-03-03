package interfaces

import "go.uber.org/zap"

type AgentActions interface {
	GetLogger() *zap.Logger
}
