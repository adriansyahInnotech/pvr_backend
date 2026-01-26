package routes

import (
	"apps/mqtt/controllers"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttRoutes struct {
	controllers *controllers.Controllers
	clientMqtt  mqtt.Client
}

func NewMqttRoutes(client mqtt.Client) *MqttRoutes {
	return &MqttRoutes{
		clientMqtt:  client,
		controllers: controllers.NewControllers(),
	}
}

func (s *MqttRoutes) Routes() {
	//subscribe all topic
	s.subscribeSafe("#", 1, s.controllers.SubscribeAllTopic)
}

func (s *MqttRoutes) subscribeSafe(topic string, qos byte, h mqtt.MessageHandler) {
	token := s.clientMqtt.Subscribe(topic, qos, h)
	token.Wait()
	if token.Error() != nil {
		log.Println("❌ Subscribe gagal:", topic, token.Error())
	} else {
		log.Println("✅ Subscribed:", topic)
	}
}
