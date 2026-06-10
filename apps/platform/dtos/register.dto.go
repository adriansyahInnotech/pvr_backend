package dtos

type Register struct {
	Username         string `json:"username"`
	Password         string `json:"password" `
	Name             string `json:"name"`
	Confirm_Password string `json:"confirm_password" `
	Secreet          string `json:"secreet"`
}
