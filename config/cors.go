package config

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Daftar origins yang diizinkan untuk CORS
var origins = []string{
	"https://auth.ulbi.ac.id",
	"https://sip.ulbi.ac.id",
	"https://euis.ulbi.ac.id",
	"https://home.ulbi.ac.id",
	"https://alpha.ulbi.ac.id",
	"https://dias.ulbi.ac.id",
	"https://iteung.ulbi.ac.id",
	"https://whatsauth.github.io",
	"https://rayfanaqbil.github.io",
	"http://127.0.0.1:5500",
	"https://hrisz.github.io",
}

// Mendapatkan host internal dari environment variables
var Internalhost string = os.Getenv("INTERNALHOST") + ":" + os.Getenv("PORT")

// Konfigurasi CORS
var Cors = cors.Config{
	AllowOrigins:     strings.Join(origins, ","),
	AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cross-Origin-Opener-Policy, Cross-Origin-Embedder-Policy",
	ExposeHeaders:    "Content-Length",
	AllowCredentials: true,
}
