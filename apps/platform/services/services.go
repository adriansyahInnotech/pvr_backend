package services

import (
	"apps/platform/dtos"
	"context"
	"encoding/json"
	"pvr_backend/config"
	"pvr_backend/helper"
	"pvr_backend/helper/response/dto"

	"github.com/gofiber/fiber/v2"
)

type PlatformServices interface {
	Publish(tracerCtx context.Context, data *dtos.PersonPublish) *dto.Response
}

type platformServices struct {
	helper     *helper.Helper
	clientMqtt *config.MqttClient
}

func NewPlatformServices(helper *helper.Helper, clientMqtt *config.MqttClient) PlatformServices {
	return &platformServices{
		clientMqtt: clientMqtt,
		helper:     helper,
	}
}

func (s *platformServices) Publish(tracerCtx context.Context, data *dtos.PersonPublish) *dto.Response {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "platform_service", "publish")
	defer span.End()

	topic := "mqtt/identify/0310741090368036"

	jsonstr, err := json.Marshal(data)
	if err != nil {
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed publish message")
	}

	if err := s.clientMqtt.Publish(topic, string(jsonstr)); err != nil {
		s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", data)
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed publish message")
	}

	return s.helper.Response.JSONResponseSuccess("", 0, 0, "berhasil mengirim pesan ke device")

}
