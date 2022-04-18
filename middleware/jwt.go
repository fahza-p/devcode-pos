package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("SECRET")),
		ErrorHandler: jwtError,
		AuthScheme:   "JWT",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
				"data":    make(map[string]string),
			})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"data":    make(map[string]string),
		})
}
