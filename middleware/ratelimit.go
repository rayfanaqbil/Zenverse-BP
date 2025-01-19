package middleware

import (
    "time"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2"
)

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,                
		Expiration: 3 * time.Minute,  
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() 
		},
	})
}