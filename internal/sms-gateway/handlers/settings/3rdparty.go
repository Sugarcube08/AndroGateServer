package settings

import (
	"fmt"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/base"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/middlewares/permissions"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/middlewares/userauth"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/devices"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/settings"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type thirdPartyControllerParams struct {
	fx.In

	DevicesSvc  *devices.Service
	SettingsSvc *settings.Service

	Validator *validator.Validate
	Logger    *zap.Logger
}

type ThirdPartyController struct {
	base.Handler

	devicesSvc  *devices.Service
	settingsSvc *settings.Service
}

func NewThirdPartyController(params thirdPartyControllerParams) *ThirdPartyController {
	return &ThirdPartyController{
		Handler: base.Handler{
			Logger:    params.Logger,
			Validator: params.Validator,
		},
		devicesSvc:  params.DevicesSvc,
		settingsSvc: params.SettingsSvc,
	}
}

// Get settings.
func (h *ThirdPartyController) get(userID string, c *fiber.Ctx) error {
	settings, err := h.settingsSvc.GetSettings(userID, true)
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	return c.JSON(settings)
}

// Update settings.
func (h *ThirdPartyController) put(userID string, c *fiber.Ctx) error {
	if err := h.BodyParserValidator(c, new(smsgateway.DeviceSettings)); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid settings format: %v", err))
	}

	settings := make(map[string]any)

	if err := c.BodyParser(&settings); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	updated, err := h.settingsSvc.ReplaceSettings(userID, settings)

	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return c.JSON(updated)
}

// Partially update settings.
func (h *ThirdPartyController) patch(userID string, c *fiber.Ctx) error {
	if err := h.BodyParserValidator(c, new(smsgateway.DeviceSettings)); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid settings format: %v", err))
	}

	settings := make(map[string]any)

	if err := c.BodyParser(&settings); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	updated, err := h.settingsSvc.UpdateSettings(userID, settings)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return c.JSON(updated)
}

func (h *ThirdPartyController) Register(app fiber.Router) {
	app.Get("", permissions.RequireScope(ScopeRead), userauth.WithUserID(h.get))
	app.Patch("", permissions.RequireScope(ScopeWrite), userauth.WithUserID(h.patch))
	app.Put("", permissions.RequireScope(ScopeWrite), userauth.WithUserID(h.put))
}
