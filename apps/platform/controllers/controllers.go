package controllers

import (
	"apps/platform/controllers/auth"
	"apps/platform/controllers/platform"
	"apps/platform/services"
	"pvr_backend/helper"
)

type Controllers struct {
	Platform *platform.PlatformControllers
	Auth     *auth.AuthController
}

func NewControllers(helper *helper.Helper, allservices *services.Services) *Controllers {

	return &Controllers{
		Platform: platform.NewPlatformControllers(helper, allservices.Platform),
		Auth:     auth.NewAuthController(helper, allservices.Auth),
	}
}
