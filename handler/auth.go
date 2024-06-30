package handler

import (
    "net/http"
    "github.com/gofiber/fiber/v2"
    "github.com/rayfanaqbil/Zenverse-BP/config"
    cek "github.com/rayfanaqbil/zenverse-BE/module"
    inimodel "github.com/rayfanaqbil/zenverse-BE/model"
)

func Login(c *fiber.Ctx) error {
    var loginDetails inimodel.Admin
    if err := c.BodyParser(&loginDetails); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Invalid request",
        })
    }

    user, err := cek.Login(config.Ulbimongoconn, "Admin", loginDetails.User_name, loginDetails.Password)
    if err != nil {
        if err.Error() == "user not found" || err.Error() == "invalid password" {
            return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
                "status":  http.StatusUnauthorized,
                "message": err.Error(),
            })
        }
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": err.Error(),
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  http.StatusOK,
        "message": "Login successful",
        "token":   user.Token,
    })
}

func Logout(c *fiber.Ctx) error {
    username := c.Locals("username").(string)

    err := cek.Logout(config.Ulbimongoconn, "Admin", username)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  fiber.StatusInternalServerError,
            "message": "Failed to logout",
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  http.StatusOK,
        "message": "Logout successful",
    })
}
