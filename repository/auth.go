// repository/guestbook_repository.go
package repository

import (
	"pvr_backend/repository/auth"

	"gorm.io/gorm"
)

type AuthRepository struct {
	User auth.UserRepository
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		User: auth.NewUserRepository(),
	}
}

// Fungsi "Sakti" untuk meng-clone seluruh repository ke mode transaksi
func (s *AuthRepository) WithTransaction(tx *gorm.DB) *AuthRepository {
	return &AuthRepository{
		User: s.User.WithTransaction(tx),
	}
}
