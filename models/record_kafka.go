package models

import (
	"time"

	"gorm.io/gorm"
)

type RecordKafka struct {
	ID        string         `gorm:"column:id;primaryKey;type:varchar(255);not null"`
	NPSN      string         `gorm:"column:npsn;type:varchar(255);index"`
	SN        string         `gorm:"column:sn;type:varchar(255);index"`
	NISN      string         `gorm:"column:nisn;type:varchar(255);index"`
	Timestamp time.Time      `gorm:"column:timestamp;type:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Area   AreaKafka   `gorm:"foreignKey:NPSN;references:NPSN;constraint:OnUpdate:CASCADE"`
	Device DeviceKafka `gorm:"foreignKey:SN;references:SN;constraint:OnUpdate:CASCADE"`
	User   UserKafka   `gorm:"foreignKey:NISN;references:NISN;constraint:OnUpdate:CASCADE"`
}

func (RecordKafka) TableName() string {
	return "record_kafka"
}
