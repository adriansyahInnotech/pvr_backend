package main

import (
	"apps/mqtt/routes"
	"log"
	"os"
	"pvr_backend/config"
)

func main() {

	client := config.NewMQTTClient(os.Getenv("RABBITMQ_MQTT_URL"), os.Getenv("RABBITMQ_MQTT_USERNAME"), os.Getenv("RABBITMQ_MQTT_PASSWORD"))
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}

	routes.NewMqttRoutes(client.Client).Routes()

	select {} // keep running
}
