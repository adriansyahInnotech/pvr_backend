package models

import (
	"time"

	"gorm.io/gorm"
)

type UserKafka struct {
	ID          uint64         `gorm:"column:id;autoIncrement;uniqueIndex"`
	NISN        string         `gorm:"column:nisn;primaryKey;type:varchar(255);not null"`
	NPSN        string         `gorm:"column:npsn;type:varchar(255);not null;index"`
	Name        string         `gorm:"column:name;type:varchar(255)"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	BiometricID *string        `gorm:"column:biometric_id;type:varchar(255)"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Area AreaKafka `gorm:"foreignKey:NPSN;references:NPSN;constraint:OnUpdate:CASCADE"`
	// SET NULL: Jika biometric dihapus, field biometric_id di User jadi NULL
	Biometric BiometricKafka `gorm:"foreignKey:BiometricID;references:ID;constraint:OnUpdate:CASCADE"`
}

func (UserKafka) TableName() string {
	return "user_kafka"
}
