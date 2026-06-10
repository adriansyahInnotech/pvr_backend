package dtos

type MQTTConnectionParams struct {
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}
