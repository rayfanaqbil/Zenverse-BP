package middleware

import (
    "github.com/gofiber/fiber/v2"
)

func VerifyCSRFToken(c *fiber.Ctx) error {
    csrfToken := c.Get("X-CSRF-Token")
    cookieCSRFToken := c.Cookies("csrf_token")
    if csrfToken == "" || csrfToken != cookieCSRFToken {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "message": "Invalid CSRF token",
        })
    }
    return nil
}

func CSRFProtection() fiber.Handler {
    return func(c *fiber.Ctx) error {
        if err := VerifyCSRFToken(c); err != nil {
            return err
        }
        return c.Next()
    }
}
