package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    iniconfig "github.com/rayfanaqbil/zenverse-BE/v2/config"
)

func Protected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("Authorization")

        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
        }

        token, err := iniconfig.ValidateJWT(tokenString)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Locals("username", claims["username"])
            return c.Next()
        }

        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
    }
}
