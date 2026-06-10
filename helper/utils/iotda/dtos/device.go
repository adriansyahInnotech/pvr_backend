package dtos

type DeviceBatchPayload struct {
	NodeID     string
	DeviceName string
	GroupID    string
}

type DeviceSuccessData struct {
	NodeID   string `json:"node_id"`
	DeviceID string `json:"device_id"`
	Secret   string `json:"secret"`
	GroupID  string `json:"group_id"`
}
