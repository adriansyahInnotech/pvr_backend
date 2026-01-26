package dto

type Response struct {
	StatusCode int    `json:"status_code"`
	Data       any    `json:"data"`
	Message    string `json:"message"`
	TotalPages int64  `json:"total_pages"`
	Page       int64  `json:"page"`
}
