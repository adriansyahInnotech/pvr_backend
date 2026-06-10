// File: apps/amqp_iotda/services/SyncStudent/SyncStudent.go
package syncstudent

import (
	"apps/amqp_iotda/dtos"
	dtos_amqp_iotda "apps/amqp_iotda/dtos"
	"context"
	"fmt"
	"log"

	"pvr_backend/helper"
	"pvr_backend/repository"
)

type SyncStudent interface {
	SyncStudent(tracerCtx context.Context, payload *dtos_amqp_iotda.PayloadSyncStudent, deviceID string, command string) error
}

type syncStudent struct {
	helper             *helper.Helper
	platformRepository repository.PlatformRepository
}

func NewSyncStudent(helper *helper.Helper, platformRepository repository.PlatformRepository) SyncStudent {
	return &syncStudent{
		helper:             helper,
		platformRepository: platformRepository,
	}
}

func (s *syncStudent) SyncStudent(tracerCtx context.Context, payload *dtos_amqp_iotda.PayloadSyncStudent, deviceID string, command string) error {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "amqp_iotda.services", "Services.SyncStudent")
	defer span.End()
	log.Println("⚙️ MASUK SERVICES SyncStudent")

	response := dtos_amqp_iotda.ResponseParasSyncStudent{}

	deviceKafka, err := s.platformRepository.DeviceKafka.GetOneByDeviceID(deviceID)
	if err != nil {

		return err
	}

	userIds, err := s.platformRepository.UserKafka.GetManyIdsByNPSN(deviceKafka.NPSN)
	if err != nil {

		return err
	}

	deviceMap := make(map[uint64]bool)
	for _, id := range payload.DeviceIDs {
		deviceMap[id] = true
	}

	var missingIDs []uint64
	for _, serverID := range userIds {
		if !deviceMap[serverID] {
			missingIDs = append(missingIDs, serverID)
		}
	}

	log.Printf("total missing id : %d", len(missingIDs))

	if len(missingIDs) == 0 {

		response = dtos_amqp_iotda.ResponseParasSyncStudent{
			Command:  command,
			Status:   "success",
			DeviceID: deviceID,
			Message:  fmt.Sprintf("✅ Alat %s sudah memiliki data terbaru, tidak ada sync diperlukan.", deviceID),
			Data:     dtos.ResponseParasDataSyncStudent{},
		}

		if err := s.helper.Utils.Iotda.SendMessageToDevice(tracerCtx, deviceID, response); err != nil {
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
			return err
		}

		log.Printf("✅ Alat %s sudah memiliki data terbaru, tidak ada sync diperlukan.", deviceID)
		return nil
	}

	missingUserKafkaModels, err := s.platformRepository.UserKafka.GetManyByIDs(missingIDs)
	if err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return err
	}

	response = dtos_amqp_iotda.ResponseParasSyncStudent{
		Command:  command,
		Status:   "success",
		DeviceID: deviceID,
		Message:  "",
		Data: dtos.ResponseParasDataSyncStudent{
			TotalData: len(*missingUserKafkaModels),
			DataSiswa: *missingUserKafkaModels,
		},
	}

	if err := s.helper.Utils.Iotda.SendMessageToDevice(tracerCtx, deviceID, response); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return err
	}

	return nil
}
