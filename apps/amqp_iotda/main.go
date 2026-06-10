package main

import (
	"context"
	"log"

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
		"d45409fb27.st1.iotda-app.ap-southeast-4.myhuaweicloud.com",
		"43xpqtlE",
		"GA1ii2qCUvG2vQTVLm1O8XE2Kmp0oVkM",
		"f0768b08-1056-4acd-8d2d-0cd468e33163",
		"DefaultQueue",
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
