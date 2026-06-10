package services

import (
	"apps/platform/services/auth"
	"apps/platform/services/platform"
	"pvr_backend/helper"
	"pvr_backend/middleware"
	"pvr_backend/repository"
)

type Services struct {
	Platform platform.PlatformServices
	Auth     auth.AuthService
}

func NewServices(helper *helper.Helper, platformRepository *repository.PlatformRepository, authRepository *repository.AuthRepository, middleware *middleware.Middleware) *Services {
	return &Services{
		Platform: platform.NewPlatformServices(helper, platformRepository),
		Auth:     auth.NewAuthService(helper, authRepository, platformRepository, middleware),
	}
}
