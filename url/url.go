package url

import (
    "github.com/rayfanaqbil/Zenverse-BP/controller"
    "github.com/rayfanaqbil/Zenverse-BP/handlers"
    "github.com/rayfanaqbil/Zenverse-BP/middleware"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/swagger"
)

func Web(page *fiber.App) {
    page.Get("/", controller.Sink)
    page.Post("/", controller.Sink)
    page.Put("/", controller.Sink)
    page.Patch("/", controller.Sink)
    page.Delete("/", controller.Sink)
    page.Options("/", controller.Sink)

    page.Get("/games", controller.GetAllGames)
    page.Get("/games/search", controller.GetGameByName)
    page.Get("/games/:id", controller.GetGamesByID)
    page.Put("/update/:id", controller.UpdateDataGames)
    page.Delete("/delete/:id", controller.DeleteGamesByID)
    page.Get("/docs/*", swagger.HandlerDefault)
    page.Post("/insert", controller.InsertDataGames)
    page.Post("/login", handlers.Login)
    page.Get("/admin", controller.GetDataAdmin)

    // Protected routes
    protected := page.Group("/api", middleware.Protected())
    protected.Get("/protected-route", func(c *fiber.Ctx) error {
        username := c.Locals("username")
        return c.JSON(fiber.Map{"message": "This is a protected route", "user": username})
    })
}
