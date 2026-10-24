package db

import (
	"github.com/jaevor/go-nanoid"
	"go.uber.org/fx"

	healthmod "github.com/Sugarcube08/AndroGateServer/pkg/health"
)

const (
	idSize = 21
)

type IDGen func() string

func Module() fx.Option {
	return fx.Module(
		"db",
		fx.Provide(
			healthmod.AsHealthProvider(newHealth),
		),
		fx.Provide(func() (IDGen, error) {
			return nanoid.Standard(idSize)
		}),
	)
}
