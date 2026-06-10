package main

import (
	"apps/platform/routes"
	"apps/platform/services"
	"context"
	"log"
	"os"
	"pvr_backend/config"
	"pvr_backend/db"
	"pvr_backend/helper"
	"pvr_backend/middleware"
	"pvr_backend/repository"
)

func init() {
	db.InitDB()
	config.InitTracer()
}

func main() {

	defer func() {
		if err := config.ShutdownTracer(context.Background()); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	// =========================
	// IOTDA SDK
	// =========================

	iotdaClient, err := config.NewIoTDAClient()

	if err != nil {
		panic(err)
	}

	log.Println("✅ IOTDA SDK CONNECTED")

	helper := helper.NewHelper(iotdaClient, nil)

	app := config.LoadConfigApp()

	middleware := middleware.NewMiddlware()

	platformRepository := repository.NewPlatformRepository()
	authRepository := repository.NewAuthRepository()

	allservices := services.NewServices(helper, platformRepository, authRepository, middleware)

	routes.NewRoutes(helper, allservices, middleware).Routes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}
