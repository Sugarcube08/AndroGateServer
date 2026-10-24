package devices

import (
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/devices"
	"github.com/Sugarcube08/AndroGateServer/internal/worker/executor"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"devices",
		logger.WithNamedLogger("devices"),
		fx.Provide(func(c Config) CleanupConfig {
			return c.Cleanup
		}, fx.Private),
		fx.Provide(devices.NewRepository, fx.Private),
		fx.Provide(
			executor.AsWorkerTask(NewCleanupTask),
		),
	)
}
