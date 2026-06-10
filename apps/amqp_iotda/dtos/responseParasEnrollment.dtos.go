package dtos

type ResponseParasEnrollmentSync struct {
	Status  string `json:"status"`
	Nisn    string `json:"nisn"`
	Message string `json:"message"`
}

type DeviceReply struct {
	ResultCode int `json:"result_code"`
}
