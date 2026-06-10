package models

import (
	"time"

	"gorm.io/gorm"
)

type AreaKafka struct {
	NPSN         string         `gorm:"column:npsn;primaryKey;type:varchar(255);not null"`
	Name         string         `gorm:"column:name;type:varchar(255)"`
	GroupIDCloud string         `gorm:"column:group_id_cloud;type:varchar;size:255"`
	CreateAt     time.Time      `gorm:"column:create_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Cascade: Jika Area dihapus, semua Device, User, dan Record yang NPSN-nya sama akan dihapus
	Devices []DeviceKafka `gorm:"foreignKey:NPSN;references:NPSN;constraint:OnUpdate:CASCADE"`
	Users   []UserKafka   `gorm:"foreignKey:NPSN;references:NPSN;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Records []RecordKafka `gorm:"foreignKey:NPSN;references:NPSN;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (AreaKafka) TableName() string {
	return "area_kafka"
}
