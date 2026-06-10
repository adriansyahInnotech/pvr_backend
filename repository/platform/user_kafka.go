package platform

import (
	"pvr_backend/db"
	"pvr_backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserKafkaRepository interface {
	WithTransaction(tx *gorm.DB) UserKafkaRepository
	UpsertBulk(userAreaModels []models.UserKafka) error
	Upsert(userAreaModels *models.UserKafka) error
	GetOneByNISN(nisn string) (*models.UserKafka, error)
	GetManyByNISN(nisn []string) (*[]models.UserKafka, error)
	GetManyIdsByNPSN(npsn string) ([]uint64, error)
	GetManyByIDs(ids []uint64) (*[]models.UserKafka, error)
}

type userKafkaRepository struct {
	db *gorm.DB
}

func NewUserKafkaRepository() UserKafkaRepository {
	return &userKafkaRepository{
		db: db.GetDB(),
	}
}

func (s *userKafkaRepository) WithTransaction(tx *gorm.DB) UserKafkaRepository {
	return &userKafkaRepository{
		db: tx,
	}
}

func (s *userKafkaRepository) UpsertBulk(userAreaModels []models.UserKafka) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "nisn"}},
		DoUpdates: clause.AssignmentColumns([]string{"npsn", "name", "biometric_id"}),
	}).CreateInBatches(&userAreaModels, 100).Error
}

func (s *userKafkaRepository) Upsert(userAreaModels *models.UserKafka) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "nisn"}},
		DoUpdates: clause.AssignmentColumns([]string{"npsn", "name", "biometric_id"}),
	}).CreateInBatches(userAreaModels, 100).Error
}

func (s *userKafkaRepository) GetOneByNISN(nisn string) (*models.UserKafka, error) {
	userKafkaModel := new(models.UserKafka)
	result := s.db.Preload("Biometric").Where("nisn = ? ", nisn).First(userKafkaModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return userKafkaModel, nil
}

func (s *userKafkaRepository) GetManyByNISN(nisn []string) (*[]models.UserKafka, error) {
	userKafkaModel := new([]models.UserKafka)
	result := s.db.Where("nisn IN ? ", nisn).Find(userKafkaModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return userKafkaModel, nil
}

func (s *userKafkaRepository) GetManyByIDs(ids []uint64) (*[]models.UserKafka, error) {
	userKafkaModel := new([]models.UserKafka)
	result := s.db.Preload("Biometric").Where("id IN ? ", ids).Find(userKafkaModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return userKafkaModel, nil
}

func (s *userKafkaRepository) GetManyIdsByNPSN(npsn string) ([]uint64, error) {
	idsUserKafka := []uint64{}
	result := s.db.Model(&models.UserKafka{}).Where("npsn = ?", npsn).Pluck("id", &idsUserKafka)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return idsUserKafka, nil
}
