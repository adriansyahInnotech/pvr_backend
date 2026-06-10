package auth

import (
	"pvr_backend/db"
	"pvr_backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	Register(user *models.User) (*models.User, error)
	WithTransaction(tx *gorm.DB) UserRepository
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: db.GetDB(),
	}
}

func (s *userRepository) WithTransaction(tx *gorm.DB) UserRepository {
	return &userRepository{
		db: tx,
	}
}

func (s *userRepository) GetUserByUsername(username string) (*models.User, error) {

	user := new(models.User)

	if err := s.db.First(&user, "username = ?", username).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return user, nil
}

func (s *userRepository) Register(user *models.User) (*models.User, error) {
	result := s.db.Create(&user)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return user, nil
}
