package utils

import (
	"pvr_backend/config"
	"pvr_backend/helper/utils/amqp"
	"pvr_backend/helper/utils/iotda"
	"pvr_backend/helper/utils/jaeger"

	iotdaSDK "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
)

type Utils struct {
	JaegerTracer jaeger.JaegerTracer
	Iotda        iotda.IotdaUtils
	Amqp         amqp.AmqpUtils
}

// File: helper/utils/utils.go

func NewUtils(
	iotdaClient *iotdaSDK.IoTDAClient,
	amqpClient *config.AMQPClient,
) *Utils {

	jaegerTracer := *jaeger.NewJaegerTracer()

	// 1. Inisialisasi Utils dasar (Tracer selalu jalan)
	u := &Utils{
		JaegerTracer: jaegerTracer,
	}

	// 2. CEK: Jika iotdaClient ADA, baru rakit IotdaUtils
	if iotdaClient != nil {
		u.Iotda = *iotda.NewIotdaUtils(iotdaClient, &jaegerTracer)
	}

	// 3. CEK: Jika amqpClient ADA, baru rakit AmqpUtils
	if amqpClient != nil {
		u.Amqp = *amqp.NewAmqpUtils(amqpClient)
	}

	return u
}
