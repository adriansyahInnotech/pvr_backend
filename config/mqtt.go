package config

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	Client MQTT.Client
}

func NewMQTTClient(broker, username, password string) *MqttClient {
	opts := MQTT.NewClientOptions().
		AddBroker(broker).
		SetUsername(username).
		SetPassword(password).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetMaxReconnectInterval(10 * time.Second).
		SetConnectionLostHandler(func(c MQTT.Client, err error) {
			fmt.Println("Connection lost:", err)
		}).
		SetOnConnectHandler(func(c MQTT.Client) {
			fmt.Println("Connected to MQTT broker")
		})

	client := MQTT.NewClient(opts)
	return &MqttClient{Client: client}
}

func (c *MqttClient) Connect() error {
	fmt.Println("🚀 Connecting to MQTT broker...")

	token := c.Client.Connect()

	// Tunggu sampai selesai atau timeout
	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("failed to connect: timeout")
	}

	// Cek error dari token
	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	fmt.Println("✅ Successfully connected to MQTT broker")
	return nil
}

func (c *MqttClient) Publish(topic, payload string) error {

	if token := c.Client.Publish(topic, 1, false, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
