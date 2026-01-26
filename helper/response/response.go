package response

import "pvr_backend/helper/response/dto"

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (s *Response) JSONResponseSuccess(data any, page int64, totalPages int64, message string) *dto.Response {

	return &dto.Response{
		StatusCode: 200,
		Data:       data,
		Message:    message,
		TotalPages: totalPages,
		Page:       page,
	}

}

func (s *Response) JSONResponseError(status_code int, message string) *dto.Response {

	return &dto.Response{
		StatusCode: status_code,
		Data:       nil,
		Message:    message,
	}

}
