package main

import (
	"context"
	"log"
	"os"

	"apps/amqp_iotda/routes"
	"apps/amqp_iotda/services"

	"pvr_backend/config"
	"pvr_backend/db"
	"pvr_backend/helper"
	"pvr_backend/repository"
)

func init() {
	db.InitDB()
	config.InitTracer()
}

func main() {

	// =========================
	// AMQP
	// =========================

	amqpClient, err := config.NewAMQPClient(
		os.Getenv("AMQP_IOTDA_HOST"),
		os.Getenv("AMQP_IOTDA_ACCESS_KEY"),
		os.Getenv("AMQP_IOTDA_ACCESS_CODE"),
		os.Getenv("AMQP_IOTDA_INSTANCE_ID"),
		os.Getenv("AMQP_IOTDA_QUEUE"),
	)

	if err != nil {
		panic(err)
	}

	// =========================
	// IOTDA SDK
	// =========================

	log.Println("✅ AMQP CONNECTED")

	// =========================
	// IOTDA SDK
	// =========================

	iotdaClient, err := config.NewIoTDAClient()

	if err != nil {
		panic(err)
	}

	log.Println("✅ IOTDA SDK CONNECTED")

	// =========================
	// HELPER
	// =========================

	appHelper := helper.NewHelper(
		iotdaClient,
		amqpClient,
	)

	// repository

	platformRepository := repository.NewPlatformRepository()

	allServices := services.NewServices(appHelper, *platformRepository)

	// =========================
	// ROUTER
	// =========================

	router := routes.NewRouter(
		appHelper,
		allServices,
	)

	// =========================
	// CONSUMER
	// =========================

	go appHelper.Utils.Amqp.Consume(context.Background(), router.Routes)

	log.Println("✅ CONSUMER STARTED")

	select {}

}
