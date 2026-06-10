package amqp

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"pvr_backend/config"
	"pvr_backend/helper/utils/amqp/dtos"
)

type MessageHandler func(
	payload dtos.DeviceMessageReport,
) error

type AmqpUtils struct {
	client *config.AMQPClient
}

func NewAmqpUtils(
	client *config.AMQPClient,
) *AmqpUtils {

	return &AmqpUtils{
		client: client,
	}
}

func (a *AmqpUtils) Consume(ctx context.Context, handler MessageHandler) {

	for {

		message, err := a.client.Receiver.Receive(ctx)

		if err != nil {
			log.Printf("❌ RECEIVE ERROR: %v", err)

			// reconnect total
			err = a.client.Reconnect()
			if err != nil {

				log.Printf("❌ RECONNECT FAILED: %v", err)

				time.Sleep(5 * time.Second)

				continue
			}

			log.Println("✅ RECONNECTED SUCCESS")

			continue
		}

		log.Println("📩 MESSAGE RECEIVED")

		log.Printf("RAW MESSAGE TYPE: %T", message.Value)
		log.Printf("RAW MESSAGE: %#v", message.Value)

		var payload dtos.DeviceMessageReport

		switch v := message.Value.(type) {

		case string:

			err = json.Unmarshal([]byte(v), &payload)

		case []byte:

			err = json.Unmarshal(v, &payload)

		default:

			log.Printf("❌ UNKNOWN MESSAGE TYPE: %T", v)

			_ = message.Reject(nil)

			continue
		}

		err = handler(payload)
		if err != nil {
			log.Printf("❌ HANDLER ERROR: %v", err)
			_ = message.Release()
			continue
		}

		err = message.Accept()
		if err != nil {
			log.Printf("❌ ACCEPT ERROR: %v", err)
			continue
		}

		log.Println("✅ MESSAGE ACCEPT")
	}
}
