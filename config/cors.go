package config

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Daftar origins yang diizinkan untuk CORS
var origins = []string{
	"http://127.0.0.1:5500",
	"http://127.0.0.1:5501",
	"http://127.0.0.1:8080",
	"https://rayfanaqbil.github.io",
	"https://hrisz.github.io",
}

// Mendapatkan host internal dari environment variables
var Internalhost string = os.Getenv("INTERNALHOST") + ":" + os.Getenv("PORT")

// Konfigurasi CORS
var Cors = cors.Config{
	AllowOrigins:     strings.Join(origins, ","),
	AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-CSRF-Token", 
	ExposeHeaders:    "Content-Length",
	AllowCredentials: true,
}
