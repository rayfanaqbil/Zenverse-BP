package handler

import(
    "net/http"
    cek "github.com/rayfanaqbil/zenverse-BE/v2/module"
    inimodel "github.com/rayfanaqbil/zenverse-BE/v2/model"
    "github.com/rayfanaqbil/Zenverse-BP/config"
    "github.com/gofiber/fiber/v2"

)



func Login(c *fiber.Ctx) error {
    var loginDetails inimodel.Admin
    if err := c.BodyParser(&loginDetails); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Invalid request",
        })
    }

    token, err := cek.Login(config.Ulbimongoconn, "Admin", loginDetails.User_name, loginDetails.Password)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": err.Error(),
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  http.StatusOK,
        "message": "Login successful",
        "token":   token,
    })
}