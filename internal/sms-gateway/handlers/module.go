package handlers

import (
	"github.com/android-sms-gateway/server/internal/sms-gateway/handlers/devices"
	"github.com/android-sms-gateway/server/internal/sms-gateway/handlers/events"
	"github.com/android-sms-gateway/server/internal/sms-gateway/handlers/messages"
	"github.com/android-sms-gateway/server/internal/sms-gateway/handlers/settings"
	"github.com/android-sms-gateway/server/internal/sms-gateway/handlers/thirdparty"
	"github.com/capcom6/go-infra-fx/http"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"handlers",
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("handlers")
		}),
		fx.Provide(
			http.AsRootHandler(newRootHandler),
			http.AsRootHandler(newAPICatalogHandler),
			http.AsApiHandler(newThirdPartyHandler),
			http.AsApiHandler(newMobileHandler),
		),
		fx.Provide(
			NewHealthHandler,
			messages.NewThirdPartyController,
			messages.NewMobileController,
			devices.NewThirdPartyController,
			settings.NewThirdPartyController,
			settings.NewMobileController,
			events.NewMobileController,
			fx.Private,
		),
		thirdparty.Module(),
	)
}
