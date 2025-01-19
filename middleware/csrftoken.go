package middleware

import (
    "github.com/gofiber/fiber/v2"
	"math/rand"
	"time"
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

func SetCSRFTokenCookie(c *fiber.Ctx, csrfToken string) {
    c.Cookie(&fiber.Cookie{
        Name:     "csrf_token",
        Value:    csrfToken,
        SameSite: "Strict", 
        Secure:   true,      
        HTTPOnly: true,      
    })
}

func GenerateCSRFToken(c *fiber.Ctx) error {
    csrfToken := generateRandomString(32)
    SetCSRFTokenCookie(c, csrfToken)
    return c.JSON(fiber.Map{
        "message": "CSRF token set successfully",
        "csrf_token": csrfToken,
    })
}

func generateRandomString(length int) string {
    randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
    charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    var token []byte
    for i := 0; i < length; i++ {
        token = append(token, charset[randomSource.Intn(len(charset))])
    }
    return string(token)

}
