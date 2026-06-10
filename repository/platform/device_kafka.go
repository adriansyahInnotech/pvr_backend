package platform

import (
	"pvr_backend/db"
	"pvr_backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceKafkaRepository interface {
	WithTransaction(tx *gorm.DB) DeviceKafkaRepository
	GetOneByAreaID(area_id string) (*models.DeviceKafka, error)
	GetManySNByAreaID(area_id string) ([]string, error)
	UpsertBulk(deviceKafkaModels []models.DeviceKafka) error
	BulkDeleteDeviceBySn(listSn []string) error
	GetDeviceDataForRestored(listSn []string) (*[]models.DeviceKafka, error)
	UpdateStatusCloud(deviceKafkaModels []models.DeviceKafka) error
	GetManyBySn(listSn []string) ([]models.DeviceKafka, error)
	GetOneByDeviceID(device_id string) (*models.DeviceKafka, error)
	GetOneBySn(sn string) (*models.DeviceKafka, error)
}

type deviceKafkaRepository struct {
	db *gorm.DB
}

func NewDeviceKafkaRepository() DeviceKafkaRepository {
	return &deviceKafkaRepository{
		db: db.GetDB(),
	}
}

func (s *deviceKafkaRepository) WithTransaction(tx *gorm.DB) DeviceKafkaRepository {
	return &deviceKafkaRepository{
		db: tx,
	}
}

func (s *deviceKafkaRepository) GetOneByAreaID(area_id string) (*models.DeviceKafka, error) {
	deviceModel := new(models.DeviceKafka)

	result := s.db.Where("npsn = ?", area_id).Find(deviceModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return deviceModel, nil
}

func (s *deviceKafkaRepository) GetOneByDeviceID(device_id string) (*models.DeviceKafka, error) {
	deviceModel := new(models.DeviceKafka)

	result := s.db.Where("device_id_on_cloud = ?", device_id).Find(deviceModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return deviceModel, nil
}

func (s *deviceKafkaRepository) GetOneBySn(sn string) (*models.DeviceKafka, error) {
	deviceModel := new(models.DeviceKafka)

	result := s.db.Where("sn = ?", sn).Find(deviceModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return deviceModel, nil
}

func (s *deviceKafkaRepository) GetManySNByAreaID(area_id string) ([]string, error) {
	serialNumbers := []string{}

	result := s.db.Model(&models.DeviceKafka{}).Where("npsn = ?", area_id).Pluck("sn", &serialNumbers)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return serialNumbers, nil
}

func (s *deviceKafkaRepository) GetManyBySn(listSn []string) ([]models.DeviceKafka, error) {
	modelsDeviceKafka := []models.DeviceKafka{}

	result := s.db.Model(&models.DeviceKafka{}).Where("sn IN  ?", listSn).Find(&modelsDeviceKafka)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return modelsDeviceKafka, nil
}

func (s *deviceKafkaRepository) GetDeviceDataForRestored(listSn []string) (*[]models.DeviceKafka, error) {
	deviceModel := new([]models.DeviceKafka)

	result := s.db.Unscoped().Model(&models.DeviceKafka{}).
		Preload("Area").
		Preload("Area.Users.Biometric").
		Preload("Records").
		Where("sn IN ? and deleted_at is not null", listSn).Find(deviceModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return deviceModel, nil
}

func (s *deviceKafkaRepository) UpsertBulk(deviceKafkaModels []models.DeviceKafka) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sn"}},
		DoUpdates: clause.AssignmentColumns([]string{"npsn", "brand", "timezone", "deleted_at"}),
	}).CreateInBatches(&deviceKafkaModels, 100).Error
}

func (s *deviceKafkaRepository) UpdateStatusCloud(deviceKafkaModels []models.DeviceKafka) error {

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sn"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_register_on_cloud", "device_id_on_cloud", "secreet"}),
	}).CreateInBatches(&deviceKafkaModels, 100).Error
}

func (s *deviceKafkaRepository) BulkDeleteDeviceBySn(listSn []string) error {

	result := s.db.
		Where("sn IN ?", listSn).
		Delete(&models.DeviceKafka{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
