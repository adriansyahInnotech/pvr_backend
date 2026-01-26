package routes

import (
	"apps/platform/controllers"
	"os"
	"pvr_backend/config"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	controllers *controllers.PlatformControllers
	// middleware  *middleware.Middleware
}

func NewRoutes(clientMqtt *config.MqttClient) *Routes {

	return &Routes{
		// controllers: controllers.NewTicketingController(helper, middleware),
		controllers: controllers.NewPlatformControllers(clientMqtt),
		// middleware:  middleware.NewMiddlware(),
	}
}

func (s *Routes) Routes(app *fiber.App) {
	user := app.Group("/api/user")
	user.Use(otelfiber.Middleware(otelfiber.WithServerName(os.Getenv("SERVICE_NAME"))))
	user.Post("/", s.controllers.Publish)

}
