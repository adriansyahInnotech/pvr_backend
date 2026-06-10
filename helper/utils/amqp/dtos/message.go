package dtos

type DeviceMessageReport struct {
	Resource    string `json:"resource"`
	Event       string `json:"event"`
	EventTime   string `json:"event_time"`
	EventTimeMS string `json:"event_time_ms"`
	RequestID   string `json:"request_id"`
	NotifyData  struct {
		Header struct {
			AppID     string `json:"app_id"`
			DeviceID  string `json:"device_id"`
			NodeID    string `json:"node_id"`
			ProductID string `json:"product_id"`
			GatewayID string `json:"gateway_id"`
		} `json:"header"`

		Body struct {
			Topic   string `json:"topic"`
			Content struct {
				Command string `json:"command"`
				Data    any    `json:"data"`
			} `json:"content"`
		} `json:"body"`
	} `json:"notify_data"`
}
