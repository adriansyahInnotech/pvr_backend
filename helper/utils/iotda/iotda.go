// File: helper/utils/iotda/iotda.go
package iotda

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"pvr_backend/helper/utils/iotda/dtos"
	"pvr_backend/helper/utils/jaeger" // Sesuaikan dengan path import Anda
	"sync"
	"time"

	iotda "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
)

type IotdaUtils struct {
	client       *iotda.IoTDAClient
	jaegerTracer *jaeger.JaegerTracer
}

func NewIotdaUtils(client *iotda.IoTDAClient, jaegerTracer *jaeger.JaegerTracer) *IotdaUtils {
	return &IotdaUtils{
		client:       client,
		jaegerTracer: jaegerTracer,
	}
}

// ==============================================================================
// 1. FUNGSI MESSAGE (DOWNSTREAM) - QoS 0
// ==============================================================================
// Digunakan untuk mengirim pesan ringan tanpa butuh kepastian alat menerima.
// Fire and forget.
func (h *IotdaUtils) SendMessageToDevice(
	ctx context.Context, deviceID string, message interface{},
) error {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "SendMessageToDevice")
	defer h.jaegerTracer.EndSpanSafe(span)

	request := &model.CreateMessageRequest{
		DeviceId: deviceID,
		Body: &model.DeviceMessageRequest{
			Message: &message,
		},
	}

	_, err := h.client.CreateMessage(request)
	if err != nil {
		h.jaegerTracer.RecordSpanError(span, err)
		log.Printf("❌ Gagal menembak Message: %v\n", err)
		return err
	}

	log.Println("📤 Message Downstream berhasil DITEMBAK.")
	return nil
}

// ==============================================================================
// 2. FUNGSI COMMAND SYNC (SINKRON) - QoS 1
// ==============================================================================
// Digunakan jika aplikasi Go Anda HARUS langsung tahu alat sukses/gagal saat itu juga.
// Aplikasi akan menahan eksekusi maksimal 20 detik menunggu balasan MQTTX.
func (h *IotdaUtils) SendCommandToDevice(ctx context.Context, deviceID string, serviceID string, commandName string, paras interface{}) (interface{}, error) {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "SendCommandToDevice")
	defer h.jaegerTracer.EndSpanSafe(span)

	request := &model.CreateCommandRequest{
		DeviceId: deviceID,
		Body: &model.DeviceCommandRequest{
			ServiceId:   &serviceID,
			CommandName: &commandName,
			Paras:       &paras,
		},
	}

	response, err := h.client.CreateCommand(request)
	if err != nil {
		h.jaegerTracer.RecordSpanError(span, err)
		log.Printf("❌ Sync Command GAGAL / TIMEOUT: %v\n", err)
		return nil, err
	}

	// Mengekstrak balasan dari alat (jika ada)
	var deviceReply interface{}
	if response != nil && response.Response != nil {
		deviceReply = response.Response

		// UBAH JADI JSON STRING AGAR TERBACA
		replyBytes, err := json.Marshal(deviceReply)
		if err == nil {
			log.Printf("✅ Sync Command SUKSES. Balasan dari alat: %s\n", string(replyBytes))
		} else {
			log.Printf("✅ Sync Command SUKSES (Tapi gagal membaca log: %v)\n", err)
		}
	}

	return deviceReply, nil

}

func (h *IotdaUtils) SendAsyncCommandToDevice(ctx context.Context, deviceID string, serviceID string, commandName string, paras interface{}) (string, error) {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "SendAsyncCommandToDevice")
	defer h.jaegerTracer.EndSpanSafe(span)

	if deviceID == "" {
		return "", fmt.Errorf("device ID tidak boleh kosong")
	}

	// Masukkan interface{} ke dalam variabel baru agar bisa diambil pointernya
	var parasPtr interface{} = paras

	// Aturan Huawei: Expire time (int32) dalam hitungan detik. 86400 = 24 jam.
	var expireTime int32 = 86400

	// 1. Susun request menggunakan AsyncDeviceCommandRequest
	request := &model.CreateAsyncCommandRequest{
		DeviceId: deviceID, // DeviceId di luar Body biasanya tidak butuh pointer
		Body: &model.AsyncDeviceCommandRequest{
			CommandName:  &commandName, // Pointer ke string
			ServiceId:    &serviceID,   // Pointer ke string
			Paras:        &parasPtr,    // Pointer ke interface{}
			ExpireTime:   &expireTime,  // Pointer ke int32
			SendStrategy: "delay",      // Pointer ke string
		},
	}

	// 2. Tembak API CreateAsyncCommand ke Huawei
	response, err := h.client.CreateAsyncCommand(request)
	if err != nil {
		log.Printf("❌ [Huawei API] Gagal membuat Async Command: %v", err)
		h.jaegerTracer.RecordSpanError(span, err)
		return "", err
	}

	// 3. Kembalikan Command ID untuk keperluan tracking
	if response.CommandId != nil {
		return *response.CommandId, nil
	}

	return "", fmt.Errorf("berhasil, tapi CommandID dari Huawei kosong")
}

// ==============================================================================
// FUNGSI SINGLE CREATE DEVICE
// ==============================================================================
func (h *IotdaUtils) CreateDevice(ctx context.Context, nodeID string, productID string, deviceName string) (*model.AddDeviceResponse, error) {
	request := &model.AddDeviceRequest{
		Body: &model.AddDevice{
			NodeId:     nodeID,
			ProductId:  productID,
			DeviceName: &deviceName,
		},
	}
	return h.client.AddDevice(request)
}

// ==============================================================================
// FUNGSI BATCH DENGAN EKSTRAKSI SECRET
// ==============================================================================
// Perhatikan: Return pertama diubah dari []string menjadi []DeviceSuccessData
func (h *IotdaUtils) BatchCreateDevices(ctx context.Context, productID string, devices []dtos.DeviceBatchPayload) ([]dtos.DeviceSuccessData, map[string]error) {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "BatchCreateDevices")
	defer h.jaegerTracer.EndSpanSafe(span)

	var wg sync.WaitGroup
	var mu sync.Mutex

	suksesList := []dtos.DeviceSuccessData{}
	gagalMap := make(map[string]error)

	maxConcurrentRequests := 10
	semaphore := make(chan struct{}, maxConcurrentRequests)

	for _, dev := range devices {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(nodeID, devName string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// 1. Tembak SDK Huawei
			res, err := h.CreateDevice(ctx, nodeID, productID, devName)

			mu.Lock()
			if err != nil {
				gagalMap[nodeID] = err
			} else if res != nil && res.DeviceId != nil && res.AuthInfo != nil && res.AuthInfo.Secret != nil {

				// 2. 🎯 INI DIA CARA MENGEKSTRAK SECRET DARI SDK HUAWEI!
				deviceID := *res.DeviceId
				secret := *res.AuthInfo.Secret

				// 4. Masukkan ke dalam daftar sukses
				suksesList = append(suksesList, dtos.DeviceSuccessData{
					NodeID:   nodeID,
					DeviceID: deviceID,
					Secret:   secret,
					GroupID:  dev.GroupID,
				})
			} else {
				gagalMap[nodeID] = fmt.Errorf("response dari cloud tidak memiliki secret yang valid")
			}
			mu.Unlock()

		}(dev.NodeID, dev.DeviceName)
	}

	wg.Wait()
	return suksesList, gagalMap
}

// ==============================================================================
// 6. FUNGSI SINGLE: DELETE DEVICE
// ==============================================================================
// deviceID: ID Perangkat (UUID) yang dihasilkan oleh Huawei Cloud.
func (h *IotdaUtils) DeleteDevice(ctx context.Context, deviceID string) error {
	request := &model.DeleteDeviceRequest{
		DeviceId: deviceID,
	}

	_, err := h.client.DeleteDevice(request)
	if err != nil {
		return err
	}
	return nil
}

// ==============================================================================
// 7. FUNGSI BATCH DELETE DEVICES DENGAN GOROUTINE & SEMAPHORE (ANTI-SPAM)
// ==============================================================================
// Mengembalikan map berisi daftar Device ID yang gagal beserta pesan error-nya
func (h *IotdaUtils) BatchDeleteDevices(ctx context.Context, deviceIDs []string) map[string]error {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "BatchDeleteDevices")
	defer h.jaegerTracer.EndSpanSafe(span)

	var wg sync.WaitGroup
	var mu sync.Mutex

	gagalMap := make(map[string]error)

	// SEMAPHORE: Membatasi maksimal 10 request API hapus berjalan bersamaan
	maxConcurrentRequests := 10
	semaphore := make(chan struct{}, maxConcurrentRequests)

	for _, id := range deviceIDs {
		wg.Add(1)

		// Ambil token semaphore
		semaphore <- struct{}{}

		go func(deviceID string) {
			defer wg.Done()

			// Lepaskan token setelah selesai
			defer func() { <-semaphore }()

			err := h.DeleteDevice(ctx, deviceID)

			// Gunakan Mutex karena map tidak aman diakses bersamaan (Thread-safe)
			if err != nil {
				mu.Lock()
				gagalMap[deviceID] = err
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()
	return gagalMap
}

// ==============================================================================
// 1. FUNGSI MEMBUAT GRUP BARU DI HUAWEI CLOUD
// ==============================================================================
// Fungsi ini dipanggil jika AreaID belum memiliki HuaweiGroupID di DB lokal.
// Nama grup di Huawei akan disamakan dengan AreaID (NPSN) Anda.
func (h *IotdaUtils) CreateGroupOnHuawei(ctx context.Context, areaID string) (string, error) {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "CreateGroupOnHuawei")
	defer h.jaegerTracer.EndSpanSafe(span)

	// Menyusun request sesuai dengan dokumentasi Huawei (AddDeviceGroupRequest)
	request := &model.AddDeviceGroupRequest{
		Body: &model.AddDeviceGroupDto{
			Name: &areaID, // Menggunakan area_id sebagai Nama Grup resmi
		},
	}

	// Eksekusi ke Huawei Cloud
	response, err := h.client.AddDeviceGroup(request)
	if err != nil {
		h.jaegerTracer.RecordSpanError(span, err)
		return "", err
	}

	// Menangkap ID acak (UUID) yang dialokasikan oleh Huawei
	if response.GroupId != nil {
		return *response.GroupId, nil
	}

	return "", fmt.Errorf("berhasil membuat grup, namun GroupID dari Huawei kosong")
}

// ==============================================================================
// FUNGSI UNTUK MENGHAPUS GROUP DI HUAWEI CLOUD
// ==============================================================================

func (h *IotdaUtils) DeleteGroupOnHuawei(ctx context.Context, groupID string) error {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "DeleteGroupOnHuawei")
	defer h.jaegerTracer.EndSpanSafe(span)

	if groupID == "" {
		return fmt.Errorf("group ID tidak boleh kosong")
	}

	// Susun request berdasarkan SDK Huawei
	request := &model.DeleteDeviceGroupRequest{
		GroupId: groupID,
	}

	// Eksekusi penghapusan ke Cloud
	_, err := h.client.DeleteDeviceGroup(request)
	if err != nil {
		h.jaegerTracer.RecordSpanError(span, err)
		log.Printf("❌ [Huawei API] Gagal menghapus grup %s: %v", groupID, err)
		return err
	}

	log.Printf("✅ Grup dengan ID %s berhasil dihapus dari Huawei Cloud", groupID)
	return nil
}

// ==============================================================================
// 2. FUNGSI MEMASUKKAN ALAT KE GRUP (PARALEL GOROUTINE)
// ==============================================================================
// Karena SDK Huawei Anda mengunci parameter DeviceId sebagai tunggal (string),
// fungsi ini akan memproses 'array' deviceIDs secara paralel lewat jalur belakang.
func (h *IotdaUtils) BatchAddDevicesToGroup(ctx context.Context, groupID string, deviceIDs []string) error {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "BatchAddDevicesToGroup")
	defer h.jaegerTracer.EndSpanSafe(span)

	if len(deviceIDs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	gagalCount := 0

	// Semaphore: Mengunci maksimal 10 request bersamaan agar tidak terkena spam limit Huawei
	semaphore := make(chan struct{}, 10)

	for _, devID := range deviceIDs {
		wg.Add(1)
		semaphore <- struct{}{} // Mengisi antrean slot

		go func(singleDeviceID string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Mengosongkan slot kembali setelah selesai

			// Menggunakan struct asli dari SDK Huawei Anda persis seperti di source file
			request := &model.CreateOrDeleteDeviceInGroupRequest{
				GroupId:  groupID,
				ActionId: "addDevice",    // Mengunci aksi untuk menambahkan perangkat
				DeviceId: singleDeviceID, // Memasukkan 1 ID perangkat tunggal
			}

			// Tembak API relasi ke Huawei Cloud
			_, err := h.client.CreateOrDeleteDeviceInGroup(request)
			if err != nil {
				log.Printf("❌ [Huawei API Error] Gagal memasukkan device %s ke grup %s: %v", singleDeviceID, groupID, err)

				mu.Lock()
				gagalCount++
				mu.Unlock()

				h.jaegerTracer.RecordSpanError(span, err)
			}
		}(devID)
	}

	// Menunggu seluruh Goroutine selesai memproses antrean alat
	wg.Wait()

	if gagalCount > 0 {
		return fmt.Errorf("%d dari %d perangkat gagal di-bind ke grup Huawei", gagalCount, len(deviceIDs))
	}

	return nil
}

// ==============================================================================
// FUNGSI MENGIRIM PESAN MASSAL LANGSUNG KE 1 GRUP TERTENTU (CARA RESMI HUAWEI)
// ==============================================================================
func (h *IotdaUtils) SendBatchMessageByGroupTask(ctx context.Context, taskName string, huaweiGroupID string, payload interface{}) (string, error) {
	_, span := h.jaegerTracer.StartSpan(ctx, "helper.iotda", "SendBatchMessageByGroupTask")
	defer h.jaegerTracer.EndSpanSafe(span)

	if huaweiGroupID == "" {
		return "", fmt.Errorf("Huawei Group ID tidak boleh kosong")
	}

	// 1. Bungkus payload ke format map
	payloadMap := map[string]interface{}{
		"message": payload,
	}

	// 2. 🪄 TRIK GOLANG: Masukkan ke variabel interface{} agar bisa dipointer
	var documentPayload interface{} = payloadMap

	// 2. 🪄 INI KUNCINYA: Gunakan TargetsFilter untuk menembak GroupID secara langsung
	targetsFilter := map[string]interface{}{
		"group_ids": []string{huaweiGroupID},
	}

	// 3. Susun request (Perhatikan: kita tidak pakai 'Targets' lagi)
	request := &model.CreateBatchTaskRequest{
		Body: &model.CreateBatchTask{
			TaskName:      taskName,
			TaskType:      "createMessages",
			TargetsFilter: targetsFilter, // <--- Filter berdasarkan Grup
			Document:      &documentPayload,
		},
	}

	// 4. Eksekusi API
	response, err := h.client.CreateBatchTask(request)
	if err != nil {
		log.Printf("❌ [Huawei API] Gagal membuat Batch Task Grup: %v", err)
		h.jaegerTracer.RecordSpanError(span, err)
		return "", err
	}

	if response.TaskId != nil {
		return *response.TaskId, nil
	}

	return "", fmt.Errorf("berhasil, tapi TaskID dari Huawei kosong")
}

func (h *IotdaUtils) GenerateMQTTParams(deviceID string, secret string) dtos.MQTTConnectionParams {
	// 1. Ambil Waktu UTC saat ini dengan format YYYYMMDDHH
	// Di Golang, "2006010215" adalah layout ajaib untuk YYYYMMDDHH
	timestamp := time.Now().UTC().Format("2006010215")

	// 2. Format Client ID Standar Huawei: {DeviceID}_0_0_{Timestamp}
	// "0_0" artinya perangkat langsung (bukan gateway) & tidak mengecek timestamp secara kaku
	clientID := fmt.Sprintf("%s_0_0_%s", deviceID, timestamp)

	// 3. Generate Password menggunakan algoritma HMAC-SHA256
	// Key = Timestamp, Message = Secret
	mac := hmac.New(sha256.New, []byte(timestamp))
	mac.Write([]byte(secret))
	password := hex.EncodeToString(mac.Sum(nil))

	// 4. Return seluruh parameter (Sesuaikan Hostname dengan milik Anda di gambar)
	return dtos.MQTTConnectionParams{
		ClientID: clientID,
		Username: deviceID, // Username sama persis dengan DeviceID
		Password: password,
		Hostname: "d45409fb27.st1.iotda-device.ap-southeast-4.myhuaweicloud.com", // URL IoTDA Anda
		Port:     8883,                                                           // 8883 untuk MQTTS (Secure), 1883 untuk MQTT (Non-Secure)
		Protocol: "MQTTS",
	}
}
