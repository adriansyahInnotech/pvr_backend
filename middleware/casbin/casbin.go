package casbin

import (
	"fmt"
	"pvr_backend/config"

	dtojwt "pvr_backend/middleware/jwt/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CasbinHandler struct {
}

func NewCasbinHandler() *CasbinHandler {
	fmt.Println("✅ Init CasbinHandler") // Tambahan untuk debug
	return &CasbinHandler{}
}

func (s *CasbinHandler) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari Fiber Context
		token, ok := c.Locals("JWT").(*jwt.Token)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: token not found",
			})
		}

		// Ambil custom claim
		claims, ok := token.Claims.(*dtojwt.CustomClaim)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: invalid claims",
			})
		}

		role := claims.Role  // -> subject
		path := c.Path()     // -> object
		method := c.Method() // -> action

		// Casbin check
		ok, err := config.GetEnforcer().Enforce(role, path, method)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error: Casbin",
			})
		}

		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}

		return c.Next()
	}
}
