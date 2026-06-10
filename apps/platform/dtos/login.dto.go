package dtos

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginDevice struct {
	Sn      string `json:"sn" validate:"required"`
	Secreet string `json:"secreet" validate:"required"`
}
