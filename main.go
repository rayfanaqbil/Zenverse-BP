package main

import (
	"log"
	"os"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/rayfanaqbil/Zenverse-BP/url"

	"github.com/gofiber/fiber/v2"
	_ "github.com/rayfanaqbil/Zenverse-BP/docs"
)

// @title TES SWAGGER Data Games
// @version 1.0
// @description This is a sample swagger for Fiber

// @contact.name API Support
// @contact.url https://github.com/rayfanaqbil
// @contact.email 714220044.@std.ulbi.ac.id

// @host zenversegames-ba223a40f69e.herokuapp.com
// @BasePath /
// @schemes https http

func main() {
	db := config.Ulbimongoconn
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))

	// Tambahkan middleware COOP dan COEP
	site.Use(config.CoopCoepMiddleware())

	site.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})
	url.Web(site, db)
	log.Fatal(site.Listen(musik.Dangdut()))
	log.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID")) // Hanya untuk debug, jangan gunakan di production!
	log.Println("GOOGLE_CLIENT_SECRET:", os.Getenv("GOOGLE_CLIENT_SECRET"))
}
