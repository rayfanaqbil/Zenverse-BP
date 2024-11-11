// url.go

package url

import (
	"github.com/rayfanaqbil/Zenverse-BP/controller"
	"github.com/rayfanaqbil/Zenverse-BP/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/mongo"
    "github.com/rayfanaqbil/Zenverse-BP/middleware"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/rayfanaqbil/Zenverse-BP/config"
)

func Web(page *fiber.App, db *mongo.Database) {

    page.Use(cors.New(config.Cors))

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
    page.Post("/login", handler.Login) 
    page.Post("/logout", handler.Logout)   
    page.Get("/auth/google/login", handler.GoogleLogin)
    page.Get("/auth/google/callback", handler.GoogleCallback)

    page.Use(middleware.AuthMiddleware())  
    page.Get("/dashboard", handler.DashboardPage)         


}
