package messages

import (
	"errors"
	"fmt"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/base"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/converters"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/handlers/middlewares/deviceauth"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/devices"
	"github.com/Sugarcube08/AndroGateServer/internal/sms-gateway/modules/messages"
	"github.com/capcom6/go-helpers/slices"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type mobileControllerParams struct {
	fx.In

	MessagesSvc *messages.Service

	Validator *validator.Validate
	Logger    *zap.Logger
}

type MobileController struct {
	base.Handler

	messagesSvc *messages.Service
}

func NewMobileController(params mobileControllerParams) *MobileController {
	return &MobileController{
		Handler: base.Handler{
			Logger:    params.Logger,
			Validator: params.Validator,
		},
		messagesSvc: params.MessagesSvc,
	}
}

// Get messages for sending.
func (h *MobileController) list(device devices.Device, c *fiber.Ctx) error {
	// Get and validate order parameter
	params := new(mobileGetQueryParams)
	if err := h.QueryParserValidator(c, params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	msgs, err := h.messagesSvc.SelectPending(device.ID, params.OrderOrDefault())
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	return c.JSON(
		smsgateway.MobileGetMessagesResponse(
			slices.Map(
				msgs,
				converters.MessageToMobileDTO,
			),
		),
	)
}

// Update message state.
func (h *MobileController) patch(device devices.Device, c *fiber.Ctx) error {
	req := smsgateway.MobilePatchMessageRequest{}
	if err := h.BodyParserValidator(c, &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	for _, v := range req {
		messageState := messages.MessageStateInput{
			ID:         v.ID,
			State:      messages.ProcessingState(v.State),
			Recipients: v.Recipients,
			States:     v.States,
		}

		err := h.messagesSvc.UpdateState(&device, messageState)
		if err != nil && !errors.Is(err, messages.ErrMessageNotFound) {
			h.Logger.Error("failed to update message status",
				zap.String("message_id", v.ID),
				zap.Error(err),
			)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MobileController) Register(router fiber.Router) {
	router.Get("", deviceauth.WithDevice(h.list))
	router.Patch("", deviceauth.WithDevice(h.patch))
}
