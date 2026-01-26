package middleware

import (
	"pvr_backend/middleware/casbin"
	"pvr_backend/middleware/jwt"
	"pvr_backend/middleware/websocket"
)

type Middleware struct {
	JWT       *jwt.Jwt
	Websocket *websocket.Websocket
	Casbin    *casbin.CasbinHandler
}

func NewMiddlware() *Middleware {
	return &Middleware{
		JWT:       jwt.NewJwt(),
		Websocket: websocket.NewWebsocket(),
		// Casbin:    casbin.NewCasbinHandler(),
	}

}
