package dtos

import "pvr_backend/models"

type ResponseParasSyncStudent struct {
	Command  string                       `json:"command"`
	Status   string                       `json:"status"`
	DeviceID string                       `json:"device_id"`
	Message  string                       `json:"message"`
	Data     ResponseParasDataSyncStudent `json:"data"`
}

type ResponseParasDataSyncStudent struct {
	TotalData int                `json:"data"`
	DataSiswa []models.UserKafka `json:"data_siswa"`
}
