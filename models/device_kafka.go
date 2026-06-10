package models

import (
	"time"

	"gorm.io/gorm"
)

type DeviceKafka struct {
	SN                string         `gorm:"column:sn;primaryKey;type:varchar(255);not null"`
	NPSN              string         `gorm:"column:npsn;type:varchar(255);not null;index"`
	CreateAt          time.Time      `gorm:"column:create_at"`
	Brand             string         `gorm:"column:brand;type:varchar(255)"`
	Timezone          int            `gorm:"column:timezone;type:int"`
	DeviceIDOnCloud   string         `gorm:"column:device_id_on_cloud;type:varchar;size:255"`
	IsRegisterOnCloud bool           `gorm:"column:is_register_on_cloud"`
	Secreet           string         `gorm:"column:secreet;type:varchar;size:255"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	Area    AreaKafka     `gorm:"foreignKey:NPSN;references:NPSN;constraint:OnUpdate:CASCADE"`
	Records []RecordKafka `gorm:"foreignKey:SN;references:SN"`
}

func (DeviceKafka) TableName() string {
	return "device_kafka"
}
