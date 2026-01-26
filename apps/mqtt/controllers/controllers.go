package controllers

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Controllers struct {
}

func NewControllers() *Controllers {
	return &Controllers{}
}

func (s *Controllers) SubscribeAllTopic(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📥 Soundbox Received [topic=%s]: \n 📥 payload : %s\n", msg.Topic(), string(msg.Payload()))
}
