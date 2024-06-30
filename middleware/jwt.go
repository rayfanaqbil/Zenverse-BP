package middleware

import (
    "context"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/dgrijalva/jwt-go"
    inimodel "github.com/rayfanaqbil/zenverse-BE/model"
    "github.com/rayfanaqbil/Zenverse-BP/config"
    "go.mongodb.org/mongo-driver/bson"
)

var jwtKey = []byte("ZnVRsERfnHRsZ")

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenStr := c.Get("Authorization")
        if tokenStr == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
        }

        tokenStr = tokenStr[len("Bearer "):]
        claims := &inimodel.Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
        }

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        db := config.Ulbimongoconn
        var user struct {
            Token string `bson:"token"`
        }
        err = db.Collection("Admin").FindOne(ctx, bson.M{"user_name": claims.UserName}).Decode(&user)
        if err != nil || user.Token != tokenStr {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
        }

        c.Locals("username", claims.UserName)
        return c.Next()
    }
}
