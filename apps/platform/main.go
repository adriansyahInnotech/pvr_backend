package main

import (
	"apps/platform/routes"
	"context"
	"log"
	"os"
	"pvr_backend/config"
	"pvr_backend/db"
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

	client := config.NewMQTTClient(os.Getenv("RABBITMQ_MQTT_URL"), os.Getenv("RABBITMQ_MQTT_USERNAME"), os.Getenv("RABBITMQ_MQTT_PASSWORD"))
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}

	app := config.LoadConfigApp()

	routes.NewRoutes(client).Routes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}
