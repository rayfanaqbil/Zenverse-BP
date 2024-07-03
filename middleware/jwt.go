package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    "strings"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Missing or malformed JWT",
            })
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid JWT format",
            })
        }

        tokenStr := parts[1]

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
            }
            return []byte("ZeNvErSERynHrSZ"), nil
        })
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid or expired JWT",
            })
        }

        // Validate token claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            if adminID, ok := claims["admin_id"].(string); ok {
                c.Locals("admin_id", adminID)
            } else {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "message": "Invalid JWT claims",
                })
            }
        } else {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid JWT",
            })
        }

        return c.Next()
    }
}
