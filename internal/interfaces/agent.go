package interfaces

import "go.uber.org/zap"

// AgentActions интерфейс для работы с действиями агент/клиента для отправки метрик в сервис
type AgentActions interface {
	GetLogger() *zap.Logger // Получение метода для логирования
}
