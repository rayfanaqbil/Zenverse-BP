package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    "github.com/rayfanaqbil/zenverse-BE/config"
    "github.com/rayfanaqbil/zenverse-BE/module"
    "go.mongodb.org/mongo-driver/mongo"
)

func Protected(db *mongo.Database) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("Authorization")

        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing JWT"})
        }

        token, err := config.ValidateJWT(tokenString)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            user_name := claims["username"].(string)
            admin, err := module.GetDataAdmin(db, user_name, "")
            if err != nil || admin.Token != tokenString {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
            }

            c.Locals("username", user_name)
            return c.Next()
        }

        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
    }
}
