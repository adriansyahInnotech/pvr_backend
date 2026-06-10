package absence

import (
	"apps/amqp_iotda/dtos"
	dtos_amqp_iotda "apps/amqp_iotda/dtos"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"pvr_backend/db"
	"pvr_backend/helper"
	"pvr_backend/repository"
	"time"

	"github.com/google/uuid"
)

type Absence interface {
	Absence(tracerCtx context.Context, payload *dtos_amqp_iotda.PayloadAbsence, deviceID string, command string) error
}

type absence struct {
	helper             *helper.Helper
	platformRepository repository.PlatformRepository
}

func NewAbsence(helper *helper.Helper, platformRepository repository.PlatformRepository) Absence {
	return &absence{
		helper:             helper,
		platformRepository: platformRepository,
	}
}

func (s *absence) Absence(tracerCtx context.Context, payload *dtos_amqp_iotda.PayloadAbsence, deviceID string, command string) error {
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "amqp_iotda.services", "Services.Enrollment")
	defer span.End()
	log.Println("⚙️ MASUK SERVICES ABSENSI")

	response := dtos_amqp_iotda.ResponseParasEnrollmentSync{}
	tx := db.GetDB().Begin()
	txRepo := s.platformRepository.WithTransaction(tx)

	if payload.IsMatching != true {

		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: "failed matching biometric in device",
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
		s.helper.Utils.JaegerTracer.RecordSpanError(span, errors.New("failed matching biometric in device"))

		log.Println("DEVICE REPLY : ", deviceReply)

		return errors.New("failed matching biometric in device")

	}

	recordsKafka, err := txRepo.RecordKafka.GetOneTodayByNisn(payload.Nisn)
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

		log.Println("DEVICE REPLY : ", deviceReply)

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err
	}

	if recordsKafka.NISN != "" {
		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("student has absence nisn : %s", payload.Nisn),
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
		s.helper.Utils.JaegerTracer.RecordSpanError(span, fmt.Errorf("student has absence nisn : %s", payload.Nisn))

		log.Println("DEVICE REPLY : ", deviceReply)

		return nil

	}

	deviceKafka, err := txRepo.DeviceKafka.GetOneByDeviceID(deviceID)
	if err != nil {

		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("device not registered on database  device_id : %s", deviceID),
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

		log.Println("DEVICE REPLY : ", deviceReply)

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err
	}

	userKafka, err := txRepo.UserKafka.GetOneByNISN(payload.Nisn)
	if err != nil {

		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("failed on server get user kafka by  nisn : %s", payload.Nisn),
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

		log.Println("DEVICE REPLY : ", deviceReply)

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err
	}

	recordsKafka.NISN = payload.Nisn
	recordsKafka.NPSN = userKafka.NPSN
	recordsKafka.SN = deviceKafka.SN
	recordsKafka.Timestamp = time.Now().UTC()
	recordsKafka.ID = uuid.NewString()

	if err := txRepo.RecordKafka.Upsert(recordsKafka); err != nil {
		response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
			Status:  "failed",
			Nisn:    payload.Nisn,
			Message: fmt.Sprintf("failed save record absence by  nisn : %s", payload.Nisn),
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

	response = dtos_amqp_iotda.ResponseParasEnrollmentSync{
		Status:  "success",
		Nisn:    payload.Nisn,
		Message: fmt.Sprintf("success absence on server with nisn : %s", payload.Nisn),
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
			Message: fmt.Sprintf("failed save record commit absence by  nisn : %s", payload.Nisn),
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

		log.Println("DEVICE REPLY : ", deviceReply)

		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, tx.Commit().Error)

		log.Println("DEVICE REPLY : ", deviceReply)

		return err
	}

	return nil
}
