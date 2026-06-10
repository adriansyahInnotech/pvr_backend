package platform

import (
	"pvr_backend/db"
	"pvr_backend/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RecordKafkaRepository interface {
	WithTransaction(tx *gorm.DB) RecordKafkaRepository
	Upsert(recordKafkaModel *models.RecordKafka) error
	GetOneTodayByNisn(nisn string) (*models.RecordKafka, error)
}

type recordKafkaRepository struct {
	db *gorm.DB
}

func NewRecordKafkaRepository() RecordKafkaRepository {
	return &recordKafkaRepository{
		db: db.GetDB(),
	}
}

func (s *recordKafkaRepository) WithTransaction(tx *gorm.DB) RecordKafkaRepository {
	return &recordKafkaRepository{
		db: tx,
	}
}

func (s *recordKafkaRepository) Upsert(recordKafkaModel *models.RecordKafka) error {

	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"npsn", "sn", "nisn", "timestamp"}),
	}).CreateInBatches(&recordKafkaModel, 100).Error

}

func (s *recordKafkaRepository) GetOneTodayByNisn(nisn string) (*models.RecordKafka, error) {

	recordKafkaModel := new(models.RecordKafka)

	now := time.Now()
	// Ambil waktu jam 00:00:00 hari ini
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	// Ambil waktu jam 23:59:59 hari ini
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.UTC().Location())

	result := s.db.Where("timestamp BETWEEN ? AND ?", startOfDay, endOfDay).Where("nisn = ? ", nisn).Find(recordKafkaModel)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return recordKafkaModel, nil

}
