package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string         `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_user_username" json:"username"`
	Name      string         `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Password  string         `gorm:"column:password;type:varchar(255);not null" json:"-"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (User) TableName() string {
	return "user"
}
