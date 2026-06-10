package platform

import (
	"pvr_backend/db"
	"pvr_backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BiometricKafkaRepository interface {
	WithTransaction(tx *gorm.DB) BiometricKafkaRepository
	Upsert(BiometricKafkaModel *models.BiometricKafka) error
}

type biometricKafkaRepository struct {
	db *gorm.DB
}

func NewBiometricKafkaRepository() BiometricKafkaRepository {
	return &biometricKafkaRepository{
		db: db.GetDB(),
	}
}

func (s *biometricKafkaRepository) WithTransaction(tx *gorm.DB) BiometricKafkaRepository {
	return &biometricKafkaRepository{
		db: tx,
	}
}

func (s *biometricKafkaRepository) Upsert(BiometricKafkaModel *models.BiometricKafka) error {

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"biometric_data_1", "biometric_data_2", "biometric_data_3", "biometric_data_4", "biometric_data_5"}),
	}).CreateInBatches(BiometricKafkaModel, 100).Error

}
