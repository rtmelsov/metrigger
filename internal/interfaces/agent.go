package interfaces

import "go.uber.org/zap"

// AgentActionsI интерфейс для работы с действиями агент/клиента для отправки метрик в сервис
type AgentActionsI interface {
	GetLogger() *zap.Logger // Получение метода для логирования
	ReportInterval() int
	PollInterval() int
	Address() string
	JwtKey() string
	RateLimit() int
}
