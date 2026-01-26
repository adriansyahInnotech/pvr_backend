package controllers

import (
	"apps/platform/dtos"
	"apps/platform/services"
	"pvr_backend/config"
	"pvr_backend/helper"

	"github.com/gofiber/fiber/v2"
)

type PlatformControllers struct {
	service services.PlatformServices
	helper  *helper.Helper
}

func NewPlatformControllers(clientMqtt *config.MqttClient) *PlatformControllers {
	return &PlatformControllers{
		service: services.NewPlatformServices(helper.NewHelper(), clientMqtt),
		helper:  helper.NewHelper(),
	}
}

func (s *PlatformControllers) Publish(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.controller", "register")
	defer span.End()

	dto := new(dtos.PersonPublish)

	if err := c.BodyParser(dto); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	data := s.service.Publish(tracectx, dto)

	return c.Status(data.StatusCode).JSON(data)

}
