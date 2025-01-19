package middleware

import (
    "context"
    "os"
    "github.com/rayfanaqbil/Zenverse-BP/config"
    "strings"
    "time"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/rayfanaqbil/zenverse-BE/v2/module"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    "google.golang.org/api/idtoken"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Missing token",
            })
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid token format",
            })
        }

        tokenString := parts[1]

        isBlacklisted, err := module.IsTokenBlacklisted(config.Ulbimongoconn, "blacklist", tokenString)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error checking token blacklist",
			})
		}
		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token has been revoked",
			})
		}
        
        if strings.Contains(tokenString, ".") { 
            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
                }
                secretKey := os.Getenv("JWT_SECRET")
                return []byte(secretKey), nil
            })
            if err != nil || !token.Valid {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "message": "Invalid JWT token",
                })
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok || !token.Valid {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "message": "Invalid token claims",
                })
            }

            adminID, ok := claims["admin_id"].(string)
            if !ok {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "message": "Invalid token claims",
                })
            }

            c.Locals("admin_id", adminID)
        } else {
            ctx := context.Background()
            payload, err := idtoken.Validate(ctx, tokenString, os.Getenv("GOOGLE_CLIENT_ID"))
            if err != nil {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "message": "Invalid Google OAuth token",
                })
            }

            googleUserID := payload.Subject 
            c.Locals("admin_id", googleUserID)
        }

        return c.Next()
    }
}

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,                
		Expiration: 3 * time.Minute,  
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() 
		},
	})
}