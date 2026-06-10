package helper

import (
	"pvr_backend/config"
	"pvr_backend/helper/response"
	"pvr_backend/helper/utils"

	iotda "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
)

type Helper struct {
	Response response.Response
	Utils    utils.Utils
}

func NewHelper(
	iotdaClient *iotda.IoTDAClient,
	amqpClient *config.AMQPClient,
) *Helper {

	return &Helper{
		Response: *response.NewResponse(),

		Utils: *utils.NewUtils(
			iotdaClient,
			amqpClient,
		),
	}
}
