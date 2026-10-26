package events

import (
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/base"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/middlewares/deviceauth"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/devices"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/sse"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MobileController struct {
	base.Handler

	sseSvc *sse.Service
}

func NewMobileController(sseService *sse.Service, validator *validator.Validate, logger *zap.Logger) *MobileController {
	return &MobileController{
		Handler: base.Handler{
			Logger:    logger,
			Validator: validator,
		},
		sseSvc: sseService,
	}
}

// Get events.
func (h *MobileController) get(device devices.Device, c *fiber.Ctx) error {
	return h.sseSvc.Handler(device.ID, c) //nolint:wrapcheck //wrapped internally
}

func (h *MobileController) Register(router fiber.Router) {
	router.Get("", deviceauth.WithDevice(h.get))
}
