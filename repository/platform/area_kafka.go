package platform

import (
	"pvr_backend/db"
	"pvr_backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AreaKafkaRepository interface {
	WithTransaction(tx *gorm.DB) AreaKafkaRepository
	Upsert(areaKafkaModel *models.AreaKafka) error
	GetByAreaID(area_id string) (*models.AreaKafka, error)
	GetManyByAreaID(area_id []string) (*[]models.AreaKafka, error)
	Delete(area_id string) error
}

type areaKafkaRepository struct {
	db *gorm.DB
}

func NewAreaKafkaRepository() AreaKafkaRepository {
	return &areaKafkaRepository{
		db: db.GetDB(),
	}
}

func (s *areaKafkaRepository) WithTransaction(tx *gorm.DB) AreaKafkaRepository {
	return &areaKafkaRepository{
		db: tx,
	}
}

func (s *areaKafkaRepository) Upsert(areaKafkaModel *models.AreaKafka) error {

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "npsn"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "group_id_cloud"}),
	}).CreateInBatches(&areaKafkaModel, 100).Error

}

func (s *areaKafkaRepository) GetByAreaID(area_id string) (*models.AreaKafka, error) {
	modelsAreaKafka := new(models.AreaKafka)

	result := s.db.Where("npsn = ?", area_id).Find(modelsAreaKafka)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return modelsAreaKafka, nil

}

func (s *areaKafkaRepository) GetManyByAreaID(area_id []string) (*[]models.AreaKafka, error) {
	modelsAreaKafka := new([]models.AreaKafka)

	result := s.db.Where("npsn IN ?", area_id).Find(modelsAreaKafka)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return modelsAreaKafka, nil

}

func (s *areaKafkaRepository) Delete(area_id string) error {

	result := s.db.
		Where("npsn = ?", area_id).
		Delete(&models.AreaKafka{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
