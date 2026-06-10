package models

type BiometricKafka struct {
	ID             string `gorm:"column:id;primaryKey;type:varchar(255);not null"`
	BiometricData1 string `gorm:"column:biometric_data_1;type:text"`
	BiometricData2 string `gorm:"column:biometric_data_2;type:text"` // Nullable (?)
	BiometricData3 string `gorm:"column:biometric_data_3;type:text"` // Nullable (?)
	BiometricData4 string `gorm:"column:biometric_data_4;type:text"` // Nullable (?)
	BiometricData5 string `gorm:"column:biometric_data_5;type:text"` // Nullable (?)
}

func (BiometricKafka) TableName() string {
	return "biometric_kafka"
}
