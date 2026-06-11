package routes

import (
	"apps/platform/controllers"
	"apps/platform/services"
	"os"
	"pvr_backend/helper"
	"pvr_backend/middleware"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	controllers *controllers.Controllers
	middleware  *middleware.Middleware
}

func NewRoutes(helper *helper.Helper, allServices *services.Services, middleware *middleware.Middleware) *Routes {

	return &Routes{
		controllers: controllers.NewControllers(helper, allServices),
		middleware:  middleware,
	}
}

func (s *Routes) Routes(app *fiber.App) {
	user := app.Group("/api/user")
	user.Use(otelfiber.Middleware(otelfiber.WithServerName(os.Getenv("SERVICE_NAME"))))

	user.Post("/login", s.controllers.Auth.Login)
	user.Post("/register", s.controllers.Auth.Register)

	user.Use(s.middleware.JWT.Handler())
	user.Post("/area/:area_id", s.controllers.Platform.AddUserArea)

	area := app.Group("/api/area")
	area.Use(otelfiber.Middleware(otelfiber.WithServerName(os.Getenv("SERVICE_NAME"))))
	area.Use(s.middleware.JWT.Handler())
	area.Post("/", s.controllers.Platform.AddArea)
	area.Patch("/:area_id", s.controllers.Platform.UpdateArea)
	area.Delete("/", s.controllers.Platform.DeleteArea)

	device := app.Group("/api/device")
	device.Post("/login", s.controllers.Auth.LoginDevice)
	device.Use(otelfiber.Middleware(otelfiber.WithServerName(os.Getenv("SERVICE_NAME"))))
	device.Use(s.middleware.JWT.Handler())
	device.Post("/", s.controllers.Platform.AddDevice)
	device.Delete("/", s.controllers.Platform.DeleteDevice)
	device.Get("/mqtt_connection", s.controllers.Platform.GetHuwaweiConnection)

}
