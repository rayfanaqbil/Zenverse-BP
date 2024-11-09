package main

import (
    "log"
    "os"

    "github.com/rayfanaqbil/Zenverse-BP/config"
    "github.com/rayfanaqbil/Zenverse-BP/url"

    "github.com/aiteung/musik"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
    _ "github.com/rayfanaqbil/Zenverse-BP/docs"
    "github.com/swaggo/fiber-swagger"
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
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to MongoDB
    db := config.Ulbimongoconn
    if db == nil {
        log.Fatal("Failed to connect to MongoDB")
    }

    // Set up Fiber app with CORS
    site := fiber.New(config.Iteung)
    site.Use(cors.New(config.Cors))

    // Middleware to pass MongoDB connection
    site.Use(func(c *fiber.Ctx) error {
        c.Locals("db", db)
        return c.Next()
    })

    // Print Google OAuth credentials for debugging
    log.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID"))
    log.Println("GOOGLE_CLIENT_SECRET:", os.Getenv("GOOGLE_CLIENT_SECRET"))

    // Route for Swagger documentation
    site.Get("/swagger/*", fiberSwagger.WrapHandler)

    // Load routes from url package
    url.Web(site, db)

    // Start the server
    log.Fatal(site.Listen(musik.Dangdut()))
}
