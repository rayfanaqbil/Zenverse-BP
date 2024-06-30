// url.go

package url

import (
	"github.com/rayfanaqbil/Zenverse-BP/controller"
	"github.com/rayfanaqbil/Zenverse-BP/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/mongo"
    "github.com/rayfanaqbil/Zenverse-BP/middleware"
)

func Web(page *fiber.App, db *mongo.Database) {
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

    page.Use(middleware.AuthMiddleware())
    page.Post("/login", handler.Login)
    page.Post("/logout", handler.Logout)
    page.Get("/token", controller.GetDataToken)
}
