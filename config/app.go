package config

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoadConfigApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("%s_%s", os.Getenv("APP_NAME"), os.Getenv("SERVICE_NAME")),
	})

	//config cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
		AllowMethods: os.Getenv("CORS_ALLOW_METHOD"),
	}))

	//healthcheck
	app.Use(healthcheck.New())

	app.Use(logger.New())

	return app
}
