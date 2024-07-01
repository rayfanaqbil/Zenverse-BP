package handler

import (
    "github.com/gofiber/fiber/v2"
    "github.com/rayfanaqbil/Zenverse-BP/config"
    cek "github.com/rayfanaqbil/zenverse-BE/v2/module"
    inimodel "github.com/rayfanaqbil/zenverse-BE/v2/model"
)

func Login(c *fiber.Ctx) error {
    var credentials inimodel.Credentials

    if err := c.BodyParser(&credentials); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse JSON",
        })
    }

    admin, token, err := cek.Login(config.Ulbimongoconn, "Admin", credentials.Username, credentials.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "login successful",
        "token":   token,
        "user":    admin,
    })
}