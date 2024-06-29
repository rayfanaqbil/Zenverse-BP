package handlers

import (
    "net/http"
    "github.com/gofiber/fiber/v2"
	cek "github.com/rayfanaqbil/zenverse-BE/module"
	inimodel "github.com/rayfanaqbil/zenverse-BE/model"
	iniconfig "github.com/rayfanaqbil/zenverse-BE/config"
	"github.com/rayfanaqbil/Zenverse-BP/config"
)

func Login(c *fiber.Ctx) error {
    var loginDetails inimodel.Admin
    if err := c.BodyParser(&loginDetails); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Invalid request",
        })
    }

    admin, err := cek.Login(config.Ulbimongoconn, "Admin", loginDetails.User_name, loginDetails.Password)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": err.Error(),
        })
    }

    token, err := iniconfig.GenerateJWT(loginDetails.User_name)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": "Failed to generate token",
        })
    }

    admin.Token = token
    if err := cek.UpdateAdminToken(config.Ulbimongoconn, admin); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": "Failed to save token",
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  http.StatusOK,
        "message": "Login successful",
        "token":   token,
    })
}