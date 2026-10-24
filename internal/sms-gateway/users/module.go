package users

import (
	"fmt"

	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/cache"
	"github.com/capcom6/go-infra-fx/db"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"users",
		logger.WithNamedLogger("users"),
		fx.Provide(func(factory cache.Factory) (*loginCache, error) {
			storage, err := factory.New("users:login")
			if err != nil {
				return nil, fmt.Errorf("can't create login cache: %w", err)
			}

			return newLoginCache(storage), nil
		}, fx.Private),
		fx.Provide(newRepository, fx.Private),
		fx.Provide(NewService),
	)
}

//nolint:gochecknoinits // framework-specific
func init() {
	db.RegisterMigration(Migrate)
}
