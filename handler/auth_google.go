package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	iniconfig "github.com/rayfanaqbil/zenverse-BE/v2/config"
	inimodel "github.com/rayfanaqbil/zenverse-BE/v2/model"
	inimodule "github.com/rayfanaqbil/zenverse-BE/v2/module"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig = &oauth2.Config{
		RedirectURL:  "https://zenversegames-ba223a40f69e.herokuapp.com/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

func GoogleLogin(c *fiber.Ctx) error {
	url := GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Code not found",
		})
	}

	token, err := GoogleOAuthConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to exchange token",
		})
	}

	client := GoogleOAuthConfig.Client(c.Context(), token) 
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to get user info",
		})
	}
	defer resp.Body.Close()

	var googleUser inimodel.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to decode user info",
		})
	}

	storedAdmin, err := inimodule.GetAdminByEmail(config.Ulbimongoconn, "Admin", googleUser.Email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "User not authorized as admin",
		})
	}


	jwtToken, err := iniconfig.GenerateJWT(*storedAdmin)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to generate token",
		})
	}

	err = inimodule.SaveTokenToDatabase(config.Ulbimongoconn, "tokens", storedAdmin.ID.Hex(), jwtToken)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to save token",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Login successful",
		"token":   jwtToken,
	})
}
