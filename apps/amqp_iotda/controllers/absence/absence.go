package absence

import (
	dtos_amqp_iotda "apps/amqp_iotda/dtos"
	"apps/amqp_iotda/services/absence"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pvr_backend/helper"
	"pvr_backend/helper/utils/amqp/dtos"
)

type Absence struct {
	serviceAbsence absence.Absence
	helper         *helper.Helper
}

func NewAbsence(serviceAbsence absence.Absence, helper *helper.Helper) *Absence {
	return &Absence{
		serviceAbsence: serviceAbsence,
		helper:         helper,
	}
}

func (s *Absence) Handle(payload dtos.DeviceMessageReport) error {

	dto := new(dtos_amqp_iotda.PayloadAbsence)

	tracectx := context.Background()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.amqp_iotda.controller", "Absence")
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

	return s.serviceAbsence.Absence(tracectx, dto, payload.NotifyData.Header.DeviceID, payload.NotifyData.Body.Content.Command)
}
