package amqprpc

import (
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"github.com/mytoolzone/task-mini-program/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, t)
	}

	return routes
}
