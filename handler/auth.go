package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"math/rand"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	iniconfig "github.com/rayfanaqbil/zenverse-BE/v2/config"
	"github.com/rayfanaqbil/zenverse-BE/v2/model"
	"github.com/rayfanaqbil/zenverse-BE/v2/module"
)

func Login(c *fiber.Ctx) error {
	var loginDetails model.Admin
	if err := c.BodyParser(&loginDetails); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
	}

	if err := validateLoginInput(loginDetails.User_name, loginDetails.Password); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	storedAdmin, err := module.GetAdminByUsername(config.Ulbimongoconn, "Admin", loginDetails.User_name)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Invalid credentials",
		})
	}


	if !iniconfig.CheckPasswordHash(loginDetails.Password, storedAdmin.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Invalid credentials",
		})
	}

	token, err := iniconfig.GenerateJWT(*storedAdmin) 
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Could not generate token",
		})
	}


	err = module.SaveTokenToDatabase(config.Ulbimongoconn, "tokens", storedAdmin.ID.Hex(), token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Could not save token",
		})
	}

	csrfToken := generateRandomString(32)
    SetCSRFTokenCookie(c, csrfToken)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Login successful",
		"token":   token,
		"csrf_token": csrfToken,
	})
}

func Logout(c *fiber.Ctx) error {
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

	token := parts[1]

	err := module.DeleteTokenFromMongoDB(config.Ulbimongoconn, "tokens", token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete token",
		})
	}

	err = module.AddToBlacklist(config.Ulbimongoconn, "blacklist", token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to blacklist token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func Register(c *fiber.Ctx) error {
	var newAdmin model.Admin
	if err := c.BodyParser(&newAdmin); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
	}

	existingAdmin, err := module.GetAdminByUsername(config.Ulbimongoconn, "Admin", newAdmin.User_name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Could not check existing username",
		})
	}
	if existingAdmin != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"status":  http.StatusConflict,
			"message": "Username already taken",
		})
	}

	hashedPassword, err := iniconfig.HashPassword(newAdmin.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Could not hash password",
		})
	}

	
	insertedID, err := module.InsertAdmin(config.Ulbimongoconn, "Admin", newAdmin.User_name, hashedPassword, newAdmin.Name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Could not register admin",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "Account registered successfully",
		"data":    insertedID,
	})
}

func DashboardPage(c *fiber.Ctx) error {
    adminID := c.Locals("admin_id")
    if adminID == nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": "Admin ID not found in context",
        })
    }

    adminIDStr := fmt.Sprintf("%v", adminID)

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  http.StatusOK,
        "message": "Dashboard access successful",
        "admin_id": adminIDStr,
    })
}	

func validateLoginInput(username, password string) error {
    re := regexp.MustCompile("^[a-zA-Z0-9_]+$")
    if !re.MatchString(username) {
        return errors.New("invalid username format")
    }

    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    
    return nil
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