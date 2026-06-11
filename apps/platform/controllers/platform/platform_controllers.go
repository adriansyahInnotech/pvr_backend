package platform

import (
	"apps/platform/dtos"
	"apps/platform/services/platform"
	"errors"
	"fmt"
	"pvr_backend/helper"
	jwtdto "pvr_backend/middleware/jwt/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PlatformControllers struct {
	service platform.PlatformServices
	helper  *helper.Helper
}

func NewPlatformControllers(helper *helper.Helper, services platform.PlatformServices) *PlatformControllers {

	return &PlatformControllers{
		service: services,
		helper:  helper,
	}
}

func (s *PlatformControllers) AddUserArea(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "UserArea")
	defer span.End()

	area_id := c.Params("area_id")

	dto := new(dtos.UserArea)

	if err := c.BodyParser(dto); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	data := s.service.AddUserArea(tracectx, dto, area_id)

	return c.Status(data.StatusCode).JSON(data)

}

func (s *PlatformControllers) AddDevice(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "AddDevice")
	defer span.End()

	dto := new(dtos.Device)

	if err := c.BodyParser(dto); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	if len(*dto) >= 100 {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, fmt.Errorf("maksimal add device per batch 100"))
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "maksimal add device per batch 100"))
	}

	data := s.service.AddDevice(tracectx, dto)

	return c.Status(data.StatusCode).JSON(data)

}

func (s *PlatformControllers) DeleteDevice(c *fiber.Ctx) error {
	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "DeleteDevice")
	defer span.End()

	dto := new([]dtos.DeleteDevice)

	if err := c.BodyParser(dto); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	if len(*dto) >= 100 {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, fmt.Errorf("maksimal delete device per batch 100"))
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "maksimal delete device per batch 100"))
	}

	data := s.service.DeleteDevice(tracectx, dto)

	return c.Status(data.StatusCode).JSON(data)
}

func (s *PlatformControllers) AddArea(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "AddArea")
	defer span.End()

	dto := new(dtos.Area)

	if err := c.BodyParser(dto); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	data := s.service.AddArea(tracectx, dto)

	return c.Status(data.StatusCode).JSON(data)

}

func (s *PlatformControllers) DeleteArea(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "DeleteArea")
	defer span.End()

	dto := new(dtos.Area)

	if err := c.BodyParser(dto); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	data := s.service.DeleteArea(tracectx, dto)

	return c.Status(data.StatusCode).JSON(data)

}

func (s *PlatformControllers) UpdateArea(c *fiber.Ctx) error {

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "UpdateArea")
	defer span.End()

	areaID := c.Params("area_id")

	if areaID == "" {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, errors.New("area id tidak boleh kosong"))
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "data tidak valid"))
	}

	dto := new(dtos.Area)

	if err := c.BodyParser(dto); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "gagal parsing body"))
	}

	dto.AreaID = areaID

	data := s.service.UpdateArea(tracectx, dto)

	return c.Status(data.StatusCode).JSON(data)

}

func (s *PlatformControllers) GetHuwaweiConnection(c *fiber.Ctx) error {
	sn := c.Locals("JWT").(*jwt.Token).Claims.(*jwtdto.CustomClaim).Sn

	if sn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "silahkan login device terlebih dahulu"))
	}

	tracectx := c.UserContext()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "AddArea")
	defer span.End()

	data := s.service.GetConnectionHuwawei(tracectx, sn)

	return c.Status(data.StatusCode).JSON(data)

}
