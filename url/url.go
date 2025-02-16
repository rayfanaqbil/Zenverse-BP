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
    page.Get("/games/apps", controller.GetAllGamesApps)
    page.Get("/games/rating", controller.GetGamesByRating)
    page.Get("/games/search", controller.GetGameByName)
    page.Get("/games/:id", controller.GetGamesByID)
    page.Get("/encrypt", controller.EncryptIDHandler)
    page.Get("/decrypt", controller.DecryptIDHandler)
    page.Put("/update/:id", controller.UpdateDataGames)
    page.Delete("/delete/:id", controller.DeleteGamesByID)
    page.Get("/docs/*", swagger.HandlerDefault)
    page.Post("/login", handler.Login) 
    page.Post("/upload/img", controller.PostUploadGithub) 
    page.Post("/logout", handler.Logout)   
    page.Post("/registeradmin", handler.Register)
    page.Post("/insert-gameadmin", controller.InsertDataGamesAdmin)
    page.Get("/csrf-token", handler.GenerateCSRFToken)
    page.Post("/insert-game",
    middleware.RateLimiter(),
    controller.InsertDataGames,)

    
    page.Put("/update-password", 
        middleware.AuthMiddleware(),
        middleware.CSRFProtection(),
        handler.UpdatePasswordAdmin) 


    page.Use(middleware.AuthMiddleware())  
    page.Get("/dashboard", handler.DashboardPage) 
}
