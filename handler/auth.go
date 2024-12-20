package handler

import (
	"github.com/rayfanaqbil/zenverse-BE/v2/module"
	"github.com/rayfanaqbil/zenverse-BE/v2/model"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	iniconfig "github.com/rayfanaqbil/zenverse-BE/v2/config"
	"github.com/gofiber/fiber/v2"
	"net/http"
    "strings"
    "fmt"
	"golang.org/x/oauth2"
	"encoding/json"
)
//login
func Login(c *fiber.Ctx) error {
	var loginDetails model.Admin
	if err := c.BodyParser(&loginDetails); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
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


	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Login successful",
		"token":   token,
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

	// Check if username already exists
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

	// Hash password
	hashedPassword, err := iniconfig.HashPassword(newAdmin.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Could not hash password",
		})
	}

	// Save admin to database
	insertedID, err := module.InsertAdmin(config.Ulbimongoconn, "Admin", newAdmin.User_name, hashedPassword, newAdmin.Email)
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

func GoogleLogin(c *fiber.Ctx) error {
	url := config.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Authorization code not provided",
		})
	}

	token, err := config.GoogleOAuthConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to exchange token",
		})
	}

	client := config.GoogleOAuthConfig.Client(c.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to get user info",
		})
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to parse user info",
		})
	}

	// Simpan atau autentikasi user berdasarkan email Google
	email := userInfo["email"].(string)
	existingAdmin, _ := module.GetAdminByEmail(config.Ulbimongoconn, "Admin", email)
	if existingAdmin == nil {
		// Daftar admin baru jika belum ada
		insertedID, err := module.InsertAdmin(config.Ulbimongoconn, "Admin", email, "", email)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "Failed to create new user",
			})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status": http.StatusOK,
			"message": "Account registered successfully",
			"data": insertedID,
		})
	}

	// Buat token JWT
	tokenString, err := iniconfig.GenerateJWT(*existingAdmin)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to generate token",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Login successful",
		"token":   tokenString,
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


