package dtos

type UserArea struct {
	Data []struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	} `json:"data"`
}
