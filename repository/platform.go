// repository/guestbook_repository.go
package repository

import (
	"pvr_backend/repository/platform"

	"gorm.io/gorm"
)

type PlatformRepository struct {
	UserKafka      platform.UserKafkaRepository
	DeviceKafka    platform.DeviceKafkaRepository
	AreaKafka      platform.AreaKafkaRepository
	BiometricKafka platform.BiometricKafkaRepository
	RecordKafka    platform.RecordKafkaRepository
}

func NewPlatformRepository() *PlatformRepository {
	return &PlatformRepository{
		UserKafka:      platform.NewUserKafkaRepository(),
		DeviceKafka:    platform.NewDeviceKafkaRepository(),
		AreaKafka:      platform.NewAreaKafkaRepository(),
		BiometricKafka: platform.NewBiometricKafkaRepository(),
		RecordKafka:    platform.NewRecordKafkaRepository(),
	}
}

// Fungsi "Sakti" untuk meng-clone seluruh repository ke mode transaksi
func (s *PlatformRepository) WithTransaction(tx *gorm.DB) *PlatformRepository {
	return &PlatformRepository{
		UserKafka:      s.UserKafka.WithTransaction(tx),
		DeviceKafka:    s.DeviceKafka.WithTransaction(tx),
		AreaKafka:      s.AreaKafka.WithTransaction(tx),
		BiometricKafka: s.BiometricKafka.WithTransaction(tx),
		RecordKafka:    s.RecordKafka.WithTransaction(tx),
	}
}
