package auth

import (
	"apps/platform/dtos"
	"apps/platform/services/auth"
	"pvr_backend/helper"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
)

type AuthController struct {
	authService auth.AuthService
	helper      helper.Helper
}

func NewAuthController(helper *helper.Helper, authService auth.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		helper:      *helper,
	}
}

func (s *AuthController) Login(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "auth-controller", "AuthController.login")
	defer span.End()

	dto := new(dtos.Login)

	if err := c.BodyParser(dto); err != nil {
		span.RecordError(err)
		span.SetStatus(fiber.StatusBadRequest, err.Error())
		span.SetAttributes(attribute.String("error", "true"))
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	data := s.authService.Login(tracectx, c, dto)

	return c.Status(data.StatusCode).JSON(data)
}

func (s *AuthController) LoginDevice(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "auth-controller", "AuthController.loginDevice")
	defer span.End()

	// PERBAIKAN: Menggunakan DTO LoginDevice
	dto := new(dtos.LoginDevice)

	if err := c.BodyParser(dto); err != nil {
		span.RecordError(err)
		span.SetStatus(fiber.StatusBadRequest, err.Error())
		span.SetAttributes(attribute.String("error", "true"))
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	data := s.authService.LoginDevice(tracectx, c, dto)

	return c.Status(data.StatusCode).JSON(data)
}

func (s *AuthController) Register(c *fiber.Ctx) error {
	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "auth-controller", "register")
	defer span.End()

	dto := new(dtos.Register)

	if err := c.BodyParser(dto); err != nil {
		span.RecordError(err)
		span.SetStatus(fiber.StatusBadRequest, err.Error())
		span.SetAttributes(attribute.String("error", "true"))
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	span.SetAttributes(
		attribute.String("request.username", dto.Username),
		attribute.String("request.ip", c.IP()),
	)

	data := s.authService.Register(tracectx, c, dto)

	return c.Status(data.StatusCode).JSON(data)
}
