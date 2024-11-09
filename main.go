package main

import (
    "log"
    "os"
    "github.com/rayfanaqbil/Zenverse-BP/config"
    "github.com/aiteung/musik"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv" 
    _ "github.com/rayfanaqbil/Zenverse-BP/docs"
    "github.com/rayfanaqbil/Zenverse-BP/url"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    db := config.Ulbimongoconn
    site := fiber.New(config.Iteung)
    site.Use(cors.New(config.Cors))
    site.Use(func(c *fiber.Ctx) error {
        c.Locals("db", db)
        return c.Next()
    })

    log.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID"))
    log.Println("GOOGLE_CLIENT_SECRET:", os.Getenv("GOOGLE_CLIENT_SECRET"))

    url.Web(site, db)
    log.Fatal(site.Listen(musik.Dangdut()))
}
