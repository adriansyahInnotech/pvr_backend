package platform

import (
	"apps/platform/dtos"
	"context"
	"fmt"
	"log"
	"pvr_backend/db"
	"pvr_backend/helper"
	"pvr_backend/helper/response/dto"
	"time"

	"pvr_backend/models"
	"pvr_backend/repository"

	iotda_helper_dtos "pvr_backend/helper/utils/iotda/dtos"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PlatformServices interface {
	AddUserArea(tracerCtx context.Context, data *dtos.UserArea, areaID string) *dto.Response
	AddArea(tracerCtx context.Context, data *dtos.Area) *dto.Response
	AddDevice(tracerCtx context.Context, data *dtos.Device) *dto.Response
	DeleteDevice(tracerCtx context.Context, data *[]dtos.DeleteDevice) *dto.Response
	GetConnectionHuwawei(tracerCtx context.Context, sn string) *dto.Response
}

type platformServices struct {
	helper             *helper.Helper
	PlatformRepository *repository.PlatformRepository
}

func NewPlatformServices(helper *helper.Helper, platformRepository *repository.PlatformRepository) PlatformServices {
	return &platformServices{
		helper:             helper,
		PlatformRepository: platformRepository,
	}
}

// area
func (s *platformServices) AddArea(tracerCtx context.Context, data *dtos.Area) *dto.Response {
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "pvr_backend.platform.services", "Services.AddArea")
	defer span.End()

	areaKafkaModel := models.AreaKafka{
		NPSN: data.AreaID,
		Name: data.Name,
	}

	if err := s.PlatformRepository.AreaKafka.Upsert(&areaKafkaModel); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed upsert")

	}

	return s.helper.Response.JSONResponseSuccess("", 0, 0, "berhasil")

}

// // user area
func (s *platformServices) AddUserArea(tracerCtx context.Context, data *dtos.UserArea, areaID string) *dto.Response {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "pvr_backend.platform.services", "Services.AddUserArea")
	defer span.End()

	tx := db.GetDB().Begin()
	txRepo := s.PlatformRepository.WithTransaction(tx)

	userKafkaModels := []models.UserKafka{}
	listNisn := []string{}
	payloadMessage := dtos.PersonPublish{
		Cmd: "student_add",
	}

	for _, v := range data.Data {
		userKafkaModel := models.UserKafka{
			NISN: v.UserID,
			NPSN: areaID,
			Name: v.Name,
		}

		listNisn = append(listNisn, v.UserID)

		userKafkaModels = append(userKafkaModels, userKafkaModel)

	}

	if err := txRepo.UserKafka.UpsertBulk(userKafkaModels); err != nil {
		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed bulk upsert")
	}

	modelsAreaKafka, err := s.PlatformRepository.AreaKafka.GetByAreaID(areaID)
	if err != nil {
		tx.Rollback()
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed get area")
	}

	for _, v := range userKafkaModels {
		infoPersons := dtos.InfoPerson{
			Nisn: v.NISN,
			Name: v.Name,
			Idx:  uint(v.ID),
		}

		payloadMessage.Data = append(payloadMessage.Data, infoPersons)
	}

	namaTask := fmt.Sprintf("student_add_%d", time.Now().Unix())

	log.Printf("🚀 Memulai broadcast ke Grup [%s] dengan Task [%s]...", modelsAreaKafka.GroupIDCloud, namaTask)

	taskID, err := s.helper.Utils.Iotda.SendBatchMessageByGroupTask(tracerCtx, namaTask, modelsAreaKafka.GroupIDCloud, payloadMessage)
	if err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "Gagal memerintahkan Huawei: "+err.Error())

	}

	if tx.Commit().Error != nil {
		tx.Rollback()
		s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", data)
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed publish message")
	}

	return s.helper.Response.JSONResponseSuccess(taskID, 0, 0, "berhasil mengirim pesan ke device")

}

// device
func (s *platformServices) AddDevice(tracerCtx context.Context, data *dtos.Device) *dto.Response {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "pvr_backend.platform.services", "Services.AddDevice")
	defer span.End()

	devicekafkaModels := []models.DeviceKafka{}
	deviceBatchUploads := []iotda_helper_dtos.DeviceBatchPayload{}
	listSn := []string{}

	// 1. MAPPING DATA
	for _, v := range *data {
		deviceKafkaModel := models.DeviceKafka{
			SN:                v.Params.Sn,
			NPSN:              v.Data.AreaID,
			Brand:             "efis",
			Timezone:          v.Data.Timezone,
			DeletedAt:         gorm.DeletedAt{},
			IsRegisterOnCloud: false,
		}

		deviceBatchUpload := iotda_helper_dtos.DeviceBatchPayload{
			NodeID:     v.Params.Sn,
			DeviceName: fmt.Sprintf("%s-%s", v.Data.AreaID, v.Params.Sn),
			GroupID:    v.Data.AreaID,
		}

		listSn = append(listSn, v.Params.Sn)
		devicekafkaModels = append(devicekafkaModels, deviceKafkaModel)
		deviceBatchUploads = append(deviceBatchUploads, deviceBatchUpload)
	}

	// 2. SIMPAN KE DATABASE LOKAL (AMAN DULUAN)
	if err := s.PlatformRepository.DeviceKafka.UpsertBulk(devicekafkaModels); err != nil {
		s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", data)
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed upsert bulk device")
	}

	// 3. PROSES CLOUD DENGAN RETRY LOGIC
	productID := "6a0ad264fc98ca6b32e609c7" // Idealnya ambil dari env (contoh: os.Getenv("HUAWEI_PRODUCT_ID"))

	devicesToRegister := deviceBatchUploads
	maxRetries := 3
	retryDelay := 2 * time.Second

	var finalSuksesList []iotda_helper_dtos.DeviceSuccessData
	var finalGagalMap map[string]error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("🔄 [Putaran %d/%d] Mendaftarkan %d perangkat ke Huawei Cloud...", attempt, maxRetries, len(devicesToRegister))

		// Eksekusi Helper (Berjalan sangat cepat & aman dari spam)
		suksesList, gagalMap := s.helper.Utils.Iotda.BatchCreateDevices(tracerCtx, productID, devicesToRegister)

		finalSuksesList = append(finalSuksesList, suksesList...)
		finalGagalMap = gagalMap

		// Jika sukses sempurna, keluar dari Loop Retry
		if len(gagalMap) == 0 {
			log.Println("✨ Sempurna! Seluruh perangkat berhasil terdaftar di Cloud.")
			break
		}

		// Jika masih ada yang gagal & masih ada sisa nyawa (retry)
		if attempt < maxRetries {
			log.Printf("⚠️ Terdapat %d perangkat gagal. Menunggu %v sebelum mencoba ulang...", len(gagalMap), retryDelay)

			// Filter HANYA perangkat yang gagal untuk putaran berikutnya
			nextRetryList := []iotda_helper_dtos.DeviceBatchPayload{}
			for _, item := range devicesToRegister {
				if _, isFailed := gagalMap[item.NodeID]; isFailed {
					nextRetryList = append(nextRetryList, item)
				}
			}
			devicesToRegister = nextRetryList
			time.Sleep(retryDelay)
		}
	}

	newDeviceKafkaModels := []models.DeviceKafka{}
	groupMapping := make(map[string][]string)
	for _, v := range finalSuksesList {

		newDeviceKafkaModel := models.DeviceKafka{
			SN:                v.NodeID,
			IsRegisterOnCloud: true,
			DeviceIDOnCloud:   v.DeviceID,
			Secreet:           v.Secret,
		}

		newDeviceKafkaModels = append(newDeviceKafkaModels, newDeviceKafkaModel)

		if v.GroupID != "" {
			groupMapping[v.GroupID] = append(groupMapping[v.GroupID], v.DeviceID)
		}

	}

	// ====================================================================================
	// BAGIAN YANG DIUBAH: LOGIKA CREATE GROUP & BINDING
	// ====================================================================================
	for areaID, devIDs := range groupMapping { // gID diganti nama jadi areaID agar tidak bingung

		// 1. CARI HUAWEI GROUP ID DARI DATABASE LOKAL
		// var huaweiGroupID string

		// TODO: Sesuaikan dengan query dan nama tabel yang Anda buat untuk menyimpan ID grup Huawei
		// Contoh:
		// var groupData models.DeviceGroup
		// errDB := s.PlatformRepository.DB.Where("area_id = ?", areaID).First(&groupData).Error
		// if errDB == nil { huaweiGroupID = groupData.HuaweiGroupID }

		modelsAreaKafka, err := s.PlatformRepository.AreaKafka.GetByAreaID(areaID)
		if err != nil {
			s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", data)
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
			return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed get data area id")
		}

		// 2. JIKA BELUM ADA DI DB, SURUH HUAWEI BUATKAN SEKARANG
		if modelsAreaKafka.GroupIDCloud == "" {
			log.Printf("📂 Grup untuk Area [%s] belum ada. Membuat grup baru di Huawei...", areaID)

			newID, errCreate := s.helper.Utils.Iotda.CreateGroupOnHuawei(tracerCtx, areaID)
			if errCreate != nil {
				log.Printf("⚠️ Gagal membuat grup Huawei untuk Area %s: %v", areaID, errCreate)
				continue // Lewati dan lanjut ke grup berikutnya
			}

			modelsAreaKafka.GroupIDCloud = newID // Gunakan ID baru ini

			// TODO: SIMPAN ID BARU INI KE DATABASE LOKAL ANDA
			// Contoh:
			// newGroup := models.DeviceGroup{ AreaID: areaID, HuaweiGroupID: newID }
			// s.PlatformRepository.DB.Create(&newGroup)
			if err := s.PlatformRepository.AreaKafka.Upsert(modelsAreaKafka); err != nil {
				s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", data)
				s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
				return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed update group id")
			}
		}

		// 3. BIND DEVICES KE HUAWEI GROUP ID YANG SUDAH VALID
		log.Printf("🔗 Memasukkan %d alat ke Grup [%s] (Huawei ID: %s)...", len(devIDs), areaID, modelsAreaKafka.GroupIDCloud)
		errGroup := s.helper.Utils.Iotda.BatchAddDevicesToGroup(tracerCtx, modelsAreaKafka.GroupIDCloud, devIDs)
		if errGroup != nil {
			log.Printf("⚠️ Gagal memasukkan sebagian alat ke grup %s: %v", areaID, errGroup)
		} else {
			log.Printf("✅ %d alat selesai dimasukkan ke grup %s", len(devIDs), areaID)
		}
	}
	// ====================================================================================
	if err := s.PlatformRepository.DeviceKafka.UpdateStatusCloud(newDeviceKafkaModels); err != nil {
		return s.helper.Response.JSONResponseError(500, "failed update status cloude register in device")
	}

	// 4. HASIL AKHIR & RESPONSE
	if len(finalGagalMap) > 0 {
		log.Printf("🚨 Selesai dengan Error: %d perangkat tetap gagal didaftarkan setelah %d percobaan.", len(finalGagalMap), maxRetries)

		responseData := map[string]interface{}{
			"saved_in_db":    len(*data),
			"cloud_success":  len(finalSuksesList),
			"cloud_failed":   len(finalGagalMap),
			"failed_details": s.formatErrorMap(finalGagalMap), // Frontend bisa tahu persis SN mana yang gagal
		}
		return s.helper.Response.JSONResponseSuccess(responseData, 0, 0, "Berhasil masuk DB lokal, namun sebagian gagal sinkron ke Cloud")
	}

	log.Println("\n\n\n\n FINAL SUKSES LIST : ", finalSuksesList)

	// Jika 100% Berhasil
	responseData := map[string]interface{}{
		"saved_in_db":   len(*data),
		"cloud_success": len(finalSuksesList),
		"cloud_failed":  0,
	}
	return s.helper.Response.JSONResponseSuccess(responseData, 0, 0, "Semua perangkat berhasil terdaftar di DB dan Cloud")
}

func (s *platformServices) DeleteDevice(tracerCtx context.Context, data *[]dtos.DeleteDevice) *dto.Response {
	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "pvr_backend.platform.services", "Services.DeleteDevice")
	defer span.End()

	listSn := []string{}

	for _, v := range *data {
		listSn = append(listSn, v.SN)
	}

	deviceKafkaModels, err := s.PlatformRepository.DeviceKafka.GetManyBySn(listSn)
	if err != nil {
		s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", data)
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed get many device by sn")
	}

	var huaweiDeviceIDs []string
	deviceIDToSNMap := make(map[string]string)

	for _, r := range deviceKafkaModels {
		if r.DeviceIDOnCloud != "" {
			huaweiDeviceIDs = append(huaweiDeviceIDs, r.DeviceIDOnCloud)
			deviceIDToSNMap[r.DeviceIDOnCloud] = r.SN
		}
	}

	// 2. PROSES PENGHAPUSAN DI CLOUD DENGAN RETRY LOGIC (Maksimal 3 Kali)
	idsToDelete := huaweiDeviceIDs
	maxRetries := 3
	retryDelay := 2 * time.Second
	var finalGagalMap map[string]error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("🔄 [Hapus Cloud - Putaran %d/%d] Menghapus %d perangkat dari Huawei Cloud...", attempt, maxRetries, len(idsToDelete))

		gagalMap := s.helper.Utils.Iotda.BatchDeleteDevices(tracerCtx, idsToDelete)
		finalGagalMap = gagalMap

		// Jika semua sukses terhapus dari Huawei Cloud, keluar dari loop retry
		if len(gagalMap) == 0 {
			log.Println("✨ Seluruh perangkat sukses dibersihkan dari Huawei Cloud.")
			break
		}

		// Jika masih ada yang gagal, filter ID yang bermasalah untuk putaran berikutnya
		if attempt < maxRetries {
			var nextRetryList []string
			for _, id := range idsToDelete {
				if _, masihGagal := gagalMap[id]; masihGagal {
					nextRetryList = append(nextRetryList, id)
				}
			}
			idsToDelete = nextRetryList
			time.Sleep(retryDelay)
		}
	}

	// 3. SINKRONISASI EVALUASI AKHIR KE DATABASE LOKAL
	var snSuksesDihapus []string
	for _, r := range deviceKafkaModels {
		// Jika DeviceID-nya tidak terdaftar di dalam finalGagalMap, artinya sukses dihapus dari Cloud
		if _, gagal := finalGagalMap[r.DeviceIDOnCloud]; !gagal {
			snSuksesDihapus = append(snSuksesDihapus, r.SN)
		}
	}

	// Eksekusi Soft Delete GORM hanya untuk perangkat yang sudah benar-benar terhapus di Cloud
	if len(snSuksesDihapus) > 0 {
		if err := s.PlatformRepository.DeviceKafka.BulkDeleteDeviceBySn(listSn); err != nil {
			s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", data)
			s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
			return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed delete bulk device")
		}
	}

	// 4. STRUKTUR RESPONSE UNTUK FRONTEND
	if len(finalGagalMap) > 0 {
		log.Printf("🚨 Selesai dengan catatan: %d perangkat gagal dihapus dari Cloud setelah %d percobaan.", len(finalGagalMap), maxRetries)

		// Konversi key error dari DeviceID (UUID) menjadi SN agar mudah dipahami oleh Frontend/User
		frontendFailedDetails := make(map[string]string)
		for devID, errCloud := range finalGagalMap {
			sn := deviceIDToSNMap[devID]
			frontendFailedDetails[sn] = errCloud.Error()
		}

		responsePartialData := map[string]interface{}{
			"total_request":   len(*data),
			"success_deleted": len(snSuksesDihapus),
			"failed_deleted":  len(finalGagalMap),
			"failed_details":  frontendFailedDetails,
		}

		return s.helper.Response.JSONResponseSuccess(responsePartialData, 0, 0, "Sebagian perangkat gagal dihapus dari cloud sehingga datanya di DB lokal dipertahankan")
	}

	// Jika sukses mutlak 100%
	responseSuccessData := map[string]interface{}{
		"total_request":   len(*data),
		"success_deleted": len(snSuksesDihapus),
		"failed_deleted":  0,
	}

	return s.helper.Response.JSONResponseSuccess(responseSuccessData, 0, 0, "Semua perangkat berhasil dihapus dari DB lokal dan Huawei Cloud")

}

func (s *platformServices) GetConnectionHuwawei(tracerCtx context.Context, sn string) *dto.Response {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(tracerCtx, "pvr_backend.platform.services", "Services.AddDevice")
	defer span.End()

	deviceKafka, err := s.PlatformRepository.DeviceKafka.GetOneBySn(sn)
	if err != nil {
		s.helper.Utils.JaegerTracer.AddObjectAsAttribute(span, "data", sn)
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(500, err.Error())
	}

	if deviceKafka.SN == "" {
		return s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "sn tidak terdaftar")
	}

	mqttConection := s.helper.Utils.Iotda.GenerateMQTTParams(deviceKafka.DeviceIDOnCloud, deviceKafka.Secreet)

	return s.helper.Response.JSONResponseSuccess(mqttConection, 0, 0, "success")
}

func (s *platformServices) formatErrorMap(errMap map[string]error) map[string]string {
	formatted := make(map[string]string)
	for sn, err := range errMap {
		formatted[sn] = err.Error()
	}
	return formatted
}
