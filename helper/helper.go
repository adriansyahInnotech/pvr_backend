package helper

import (
	"pvr_backend/helper/response"
	"pvr_backend/helper/utils"
)

type Helper struct {
	Response response.Response
	Utils    utils.Utils
}

func NewHelper() *Helper {

	return &Helper{
		Response: *response.NewResponse(),
		Utils:    *utils.NewUtils(),
	}

}
