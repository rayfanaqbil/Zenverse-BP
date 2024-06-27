package url

import (
	"github.com/rayfanaqbil/Zenverse-BP/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Web(page *fiber.App) {
	// page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)  //API from user whatsapp message from iteung gowa
	// page.Get("/ws/whatsauth/qr", websocket.New(controller.WsWhatsAuthQR)) //websocket whatsauth

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
	page.Post("/login", controller.Login)
	page.Get("/admin", controller.GetDataAdmin)
}
