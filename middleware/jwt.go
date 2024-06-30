package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/dgrijalva/jwt-go"
    "github.com/rayfanaqbil/zenverse-BE/model"
)

var jwtKey = []byte("ZnVRsERfnHRsZ")

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenStr := c.Get("Authorization")
        if tokenStr == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
        }

        // Remove the "Bearer " prefix
        tokenStr = tokenStr[len("Bearer "):]

        claims := &model.Admin{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
        }

        c.Locals("username", claims.User_name)
        return c.Next()
    }
}
