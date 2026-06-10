package syncstudent

import (
	dtos_amqp_iotda "apps/amqp_iotda/dtos"
	syncstudent "apps/amqp_iotda/services/sync_student"
	"context"
	"encoding/json"
	"log"
	"pvr_backend/helper"
	"pvr_backend/helper/utils/amqp/dtos"
)

type SyncStudent struct {
	serviceSyncStudent syncstudent.SyncStudent
	helper             *helper.Helper
}

func NewSyncStudent(serviceSyncStudent syncstudent.SyncStudent, helper *helper.Helper) *SyncStudent {
	return &SyncStudent{
		serviceSyncStudent: serviceSyncStudent,
		helper:             helper,
	}
}

func (s *SyncStudent) Handle(payload dtos.DeviceMessageReport) error {

	dto := new(dtos_amqp_iotda.PayloadSyncStudent)

	tracectx := context.Background()
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracectx, "pvr_backend.amqp_iotda.controller", "SyncStudent")
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

	return s.serviceSyncStudent.SyncStudent(tracectx, dto, payload.NotifyData.Header.DeviceID, payload.NotifyData.Body.Content.Command)
}
