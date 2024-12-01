package config

import (
	"log"
	"github.com/gofiber/fiber/v2"
)

func CoopCoepMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cross-Origin-Opener-Policy", "same-origin")
		c.Set("Cross-Origin-Embedder-Policy", "require-corp")
		log.Println("COOP/COEP middleware executed")
		return c.Next()
	}
}
