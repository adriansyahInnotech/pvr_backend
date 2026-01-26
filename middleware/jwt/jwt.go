package jwt

import (
	"fmt"
	"pvr_backend/middleware/jwt/dto"

	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
}

func NewJwt() *Jwt {
	return &Jwt{}
}

func (s *Jwt) Handler() fiber.Handler {
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECREET"))},
		ContextKey: "JWT",
		Claims:     new(dto.CustomClaim),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			fmt.Println("JWT Error:", err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	})

	return jwtMiddleware

	// return jwt
}

func (s *Jwt) CreateTokenJwt(claims *dto.CustomClaim) string {

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    fmt.Sprintf("%s-%s", os.Getenv("APP_NAME"), os.Getenv("SERVICE_NAME")),
		Subject:   "user-token",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECREET")))
	if err != nil {
		fmt.Println("failed generate token jwt")
		return ""
	}

	return t

}
