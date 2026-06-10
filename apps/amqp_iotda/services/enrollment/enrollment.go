// File: apps/amqp_iotda/services/enrollment/enrollment.go
package enrollment

import (
	"apps/amqp_iotda/dtos"
	dtos_amqp_iotda "apps/amqp_iotda/dtos"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"pvr_backend/db"
	"pvr_backend/helper"
	"pvr_backend/models"
	"pvr_backend/repository"

	"github.com/google/uuid"
)

type Enrollment interface {
	Enrollment(tracerCtx context.Context, payload *dtos_amqp_iotda.PayloadEnrollment, deviceID string, command string) error
}

type enrollment struct {
	helper             *helper.Helper
	platformRepository repository.PlatformRepository
}

func NewEnrollment(helper *helper.Helper, platformRepository repository.PlatformRepository) Enrollment {
	return &enrollment{
		helper:             helper,
		platformRepository: platformRepository,
	}
}

func (s *enrollment) Enrollment(tracerCtx context.Context, payload *dtos_amqp_iotda.PayloadEnrollment, deviceID string, command string) error {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "amqp_iotda.services", "Services.Enrollment")
	defer span.End()
	log.Println("⚙️ MASUK SERVICES ENROLLMENT")

	tx := db.GetDB().Begin()
	txRepo := s.platformRepository.WithTransaction(tx)

	response := dtos_amqp_iotda.ResponseParasEnrollmentSync{}

	biometricKafkaModel := models.BiometricKafka{}

	userKafkaModel, err := s.platformRepository.UserKafka.GetOneByNISN(payload.Nisn)
	if err != nil {

		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("failed on server searching  nisn : %s", payload.Nisn),
		}

		deviceReply, err2 := s.helper.Utils.Iotda.SendCommandToDevice(
			context.Background(),
			deviceID,
			command,
			command,
			response,
		)

		if err2 != nil {
			tx.Rollback()
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err2)
			log.Printf("❌ GAGAL MENGIRIM COMMAND BALASAN: %v\n", err2)
			return err2
		}

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err
	}

	if userKafkaModel.NISN == "" {

		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("nisn not found on server nisn : %s", payload.Nisn),
		}

		deviceReply, err2 := s.helper.Utils.Iotda.SendCommandToDevice(
			context.Background(),
			deviceID,
			command,
			command,
			response,
		)

		if err2 != nil {
			tx.Rollback()
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err2)
			log.Printf("❌ GAGAL MENGIRIM COMMAND BALASAN: %v\n", err)
			return err2
		}

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, fmt.Errorf("nisn not found on server nisn : %s", payload.Nisn))

		log.Println("DEVICE REPLY : ", deviceReply)

		return fmt.Errorf("murid tidak di temukan dengan nisn %s", payload.Nisn)
	}

	if userKafkaModel.BiometricID != nil {
		biometricKafkaModel.ID = *userKafkaModel.BiometricID
	} else {
		biometricKafkaModel.ID = uuid.NewString()
		userKafkaModel.BiometricID = &biometricKafkaModel.ID
	}

	if payload.PalmValueLeft != "" {
		biometricKafkaModel.BiometricData1 = payload.PalmValueLeft
	}

	if payload.PalmValueRight != "" {
		biometricKafkaModel.BiometricData2 = payload.PalmValueRight
	}

	if err := txRepo.BiometricKafka.Upsert(&biometricKafkaModel); err != nil {

		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("failed save enrollment biometric : %s", payload.Nisn),
		}

		deviceReply, err2 := s.helper.Utils.Iotda.SendCommandToDevice(
			context.Background(),
			deviceID,
			command,
			command,
			response,
		)

		if err2 != nil {
			tx.Rollback()
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err2)
			log.Printf("❌ GAGAL MENGIRIM COMMAND BALASAN: %v\n", err)
			return err2
		}

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err
	}

	if err := txRepo.UserKafka.Upsert(userKafkaModel); err != nil {

		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("failed save enrollment user : %s", payload.Nisn),
		}

		deviceReply, err2 := s.helper.Utils.Iotda.SendCommandToDevice(
			context.Background(),
			deviceID,
			command,
			command,
			response,
		)

		if err2 != nil {
			tx.Rollback()
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err2)
			log.Printf("❌ GAGAL MENGIRIM COMMAND BALASAN: %v\n", err)
			return err2
		}

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err

	}

	response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
		Status:  "success",
		Nisn:    payload.Nisn,
		Message: fmt.Sprintf("success enroll on server with nisn : %s", payload.Nisn),
	}

	deviceReply, err := s.helper.Utils.Iotda.SendCommandToDevice(
		context.Background(),
		deviceID,
		command,
		command,
		response,
	)

	if err != nil {
		log.Printf("❌ GAGAL MENGIRIM COMMAND BALASAN: %v\n", err)
		return err
	}

	if deviceReply == nil {

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, fmt.Errorf("device tidak membalas error"))

		return fmt.Errorf("device tidak membalas command")
	} else {
		replyByte, err := json.Marshal(deviceReply)
		if err != nil {
			log.Println("failed marshal device reply")
		}

		deviceReplyDto := new(dtos.DeviceReply)

		if err := json.Unmarshal(replyByte, deviceReplyDto); err != nil {
			log.Println("failed unmarshal device Reply")
		}

		log.Println("isi data device reply dto : ", deviceReplyDto)

		if deviceReplyDto.ResultCode == 1 {
			log.Println("device juga gagalmenyimpan di memory nya")
			tx.Rollback()
			s.helper.Utils.JaegerTracer.RecordSpanError(span, fmt.Errorf("device juga gagal menyimpan enrollment"))

			return fmt.Errorf("device juga gagal menyimpan enrolment")
		}
	}

	if tx.Commit().Error != nil {
		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("failed save enrollment user : %s", payload.Nisn),
		}

		deviceReply, err2 := s.helper.Utils.Iotda.SendCommandToDevice(
			context.Background(),
			deviceID,
			command,
			command,
			response,
		)

		if err2 != nil {
			tx.Rollback()
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err2)
			log.Printf("❌ GAGAL MENGIRIM COMMAND BALASAN: %v\n", err)
			return err2
		}

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err
	}

	return nil
}
