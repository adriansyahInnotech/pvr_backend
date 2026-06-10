package dtos

type PayloadEnrollment struct {
	DeviceID       string `json:"device_id"`
	AreaID         string `json:"area_id"`
	Nisn           string `json:"nisn"`
	PalmValueRight string `json:"palm_value_right"`
	PalmValueLeft  string `json:"palm_value_left"`
}
