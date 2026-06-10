package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaim struct {
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	Sn       string    `json:"sn"`
	Platform string    `json:"platform"`
	Name     string    `json:"name"`
	jwt.RegisteredClaims
}
