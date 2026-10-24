package tokens

import (
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/jwt"
	"github.com/Sugarcube08/AndroGateServer/internal/worker/executor"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"tokens",
		logger.WithNamedLogger("tokens"),
		fx.Provide(func(c Config) CleanupConfig {
			return c.Cleanup
		}, fx.Private),
		fx.Provide(jwt.NewRepository, fx.Private),
		fx.Provide(
			executor.AsWorkerTask(NewCleanupTask),
		),
	)
}
