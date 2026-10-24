package messages

import (
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/messages"
	"github.com/Sugarcube08/AndroGateServer/internal/worker/executor"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"messages",
		logger.WithNamedLogger("messages"),
		fx.Provide(func(c Config) (HashingConfig, CleanupConfig) {
			return c.Hashing, c.Cleanup
		}, fx.Private),
		fx.Provide(messages.NewRepository, fx.Private),
		fx.Provide(
			executor.AsWorkerTask(NewInitialHashingTask),
			executor.AsWorkerTask(NewCleanupTask),
		),
	)
}
