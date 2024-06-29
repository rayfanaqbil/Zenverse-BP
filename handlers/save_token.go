package handlers

import (
    "net/http"
    "github.com/gofiber/fiber/v2"
    cek "github.com/rayfanaqbil/zenverse-BE/module"
)

func SaveToken(c *fiber.Ctx) error {
    type TokenRequest struct {
        Token string `json:"token"`
    }

    var tokenReq TokenRequest
    if err := c.BodyParser(&tokenReq); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "status":  http.StatusBadRequest,
            "message": "Invalid request",
        })
    }

    if err := cek.SaveTokenToDatabase(tokenReq.Token); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "status":  http.StatusInternalServerError,
            "message": "Failed to save token",
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "status":  http.StatusOK,
        "message": "Token saved successfully",
    })
}
