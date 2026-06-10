package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"pack.ag/amqp"
)

type AMQPClient struct {
	Host string

	AccessKey string

	AccessCode string

	InstanceID string

	QueueName string

	Client *amqp.Client

	Session *amqp.Session

	Receiver *amqp.Receiver
}

func NewAMQPClient(
	host string,

	accessKey string,

	accessCode string,

	instanceID string,

	queueName string,
) (*AMQPClient, error) {

	amqpClient := &AMQPClient{
		Host: host,

		AccessKey: accessKey,

		AccessCode: accessCode,

		InstanceID: instanceID,

		QueueName: queueName,
	}

	err := amqpClient.Connect()

	if err != nil {
		return nil, err
	}

	return amqpClient, nil
}

func (a *AMQPClient) Connect() error {

	address := fmt.Sprintf(
		"amqps://%s:5671",
		a.Host,
	)

	username := fmt.Sprintf(
		"accessKey=%s|timestamp=%d|instanceId=%s",
		a.AccessKey,
		time.Now().UnixNano()/1000000,
		a.InstanceID,
	)

	client, err := amqp.Dial(
		address,

		amqp.ConnSASLPlain(
			username,
			a.AccessCode,
		),

		amqp.ConnProperty(
			"vhost",
			"default",
		),

		amqp.ConnServerHostname(
			"default",
		),

		amqp.ConnTLSConfig(
			&tls.Config{
				InsecureSkipVerify: true,

				MaxVersion: tls.VersionTLS12,
			},
		),

		amqp.ConnConnectTimeout(
			8*time.Second,
		),
	)

	if err != nil {
		return err
	}

	session, err := client.NewSession()

	if err != nil {

		safeClose(
			"client",

			func() {

				_ = client.Close()
			},
		)

		return err
	}

	receiver, err := session.NewReceiver(

		amqp.LinkTargetDurability(
			amqp.DurabilityUnsettledState,
		),

		amqp.LinkSourceAddress(
			a.QueueName,
		),

		amqp.LinkCredit(100),
	)

	if err != nil {

		safeClose(
			"session",

			func() {

				ctx, cancel := context.WithTimeout(
					context.Background(),
					5*time.Second,
				)

				defer cancel()

				_ = session.Close(ctx)
			},
		)

		safeClose(
			"client",

			func() {

				_ = client.Close()
			},
		)

		return err
	}

	a.Client = client

	a.Session = session

	a.Receiver = receiver

	log.Println(
		"✅ AMQP CONNECTED",
	)

	return nil
}

func (a *AMQPClient) Reconnect() error {

	log.Println(
		"🔄 RECONNECTING AMQP...",
	)

	safeClose(
		"receiver",

		func() {

			if a.Receiver != nil {

				ctx, cancel := context.WithTimeout(
					context.Background(),
					5*time.Second,
				)

				defer cancel()

				_ = a.Receiver.Close(ctx)
			}
		},
	)

	safeClose(
		"session",

		func() {

			if a.Session != nil {

				ctx, cancel := context.WithTimeout(
					context.Background(),
					5*time.Second,
				)

				defer cancel()

				_ = a.Session.Close(ctx)
			}
		},
	)

	safeClose(
		"client",

		func() {

			if a.Client != nil {

				_ = a.Client.Close()
			}
		},
	)

	a.Client = nil
	a.Session = nil
	a.Receiver = nil

	time.Sleep(
		5 * time.Second,
	)

	err := a.Connect()

	if err != nil {

		log.Printf(
			"❌ RECONNECT FAILED: %v",
			err,
		)

		return err
	}

	log.Println(
		"✅ RECONNECTED SUCCESS",
	)

	return nil
}

func safeClose(
	name string,

	fn func(),
) {

	defer func() {

		if r := recover(); r != nil {

			log.Printf(
				"⚠️ SAFE CLOSE %s PANIC: %v",
				name,
				r,
			)
		}
	}()

	fn()
}
