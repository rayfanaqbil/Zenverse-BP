package middleware

import (
    "github.com/gofiber/fiber/v2"
     "github.com/gofiber/fiber/v2/middleware/csrf"
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
    return csrf.New(csrf.Config{
        KeyLookup:      "header:X-CSRF-Token",
        CookieName:     "csrf_token",
        CookieSameSite: "Strict",
    })
}
