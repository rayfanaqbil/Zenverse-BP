package middleware

import (
    "github.com/gofiber/fiber/v2"
)

func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Set("Content-Security-Policy", "default-src 'self' https://zenversegames-ba223a40f69e.herokuapp.com; script-src 'self' https://cdn.jsdelivr.net; style-src 'self' https://cdn.jsdelivr.net;")
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-XSS-Protection", "1; mode=block")
        return c.Next()
    }
}
