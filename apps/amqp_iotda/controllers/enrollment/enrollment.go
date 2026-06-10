package enrollment

import (
	dtos_amqp_iotda "apps/amqp_iotda/dtos"
	"apps/amqp_iotda/services/enrollment"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pvr_backend/helper"
	"pvr_backend/helper/utils/amqp/dtos"
)

type Enrollment struct {
	serviceEnrollment enrollment.Enrollment
	helper            *helper.Helper
}

func NewEnrollment(serviceEnrollment enrollment.Enrollment, helper *helper.Helper) *Enrollment {
	return &Enrollment{
		serviceEnrollment: serviceEnrollment,
		helper:            helper,
	}
}

func (s *Enrollment) Handle(payload dtos.DeviceMessageReport) error {

	dto := new(dtos_amqp_iotda.PayloadEnrollment)

	tracectx := context.Background()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.platform.controller", "AddDevice")
	defer span.End()

	datastr, err := json.Marshal(payload.NotifyData.Body.Content.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(datastr, dto); err != nil {
		return err
	}

	log.Println(dto)
	log.Println(string(datastr))

	if dto.Nisn == "" {
		return fmt.Errorf("nisn required")
	}

	return s.serviceEnrollment.Enrollment(tracectx, dto, payload.NotifyData.Header.DeviceID, payload.NotifyData.Body.Content.Command)
}
