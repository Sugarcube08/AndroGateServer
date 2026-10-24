package tasks

import (
	"github.com/Sugarcube08/AndroGateServer/internal/worker/tasks/devices"
	"github.com/Sugarcube08/AndroGateServer/internal/worker/tasks/messages"
	"github.com/Sugarcube08/AndroGateServer/internal/worker/tasks/tokens"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"tasks",
		logger.WithNamedLogger("tasks"),
		messages.Module(),
		devices.Module(),
		tokens.Module(),
	)
}
