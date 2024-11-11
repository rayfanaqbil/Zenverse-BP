package handler

import (
	"context"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	iniconfig "github.com/rayfanaqbil/zenverse-BE/v2/config"
	inimodel "github.com/rayfanaqbil/zenverse-BE/v2/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

var (
	GoogleOAuthConfig = &oauth2.Config{
		RedirectURL:  "https://zenversegames-ba223a40f69e.herokuapp.com/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	allowedAdmins = []string{"rayfana09@gmail.com", "harissaefuloh@gmail.com"}
)

func GoogleLogin(c *fiber.Ctx) error {
	url := GoogleOAuthConfig.AuthCodeURL("state-token")
	return c.Redirect(url)
}

// GoogleCallback menangani callback dari Google OAuth
func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).SendString("Kode tidak ditemukan")
	}

	token, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal menukar token")
	}

	client := GoogleOAuthConfig.Client(context.Background(), token)

	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membuat layanan OAuth2")
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal mengambil info pengguna")
	}

	// Cek apakah email pengguna termasuk dalam daftar admin yang diizinkan
	isAllowed := false
	for _, adminEmail := range allowedAdmins {
		if userInfo.Email == adminEmail {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return c.Status(fiber.StatusForbidden).SendString("Akses hanya untuk admin")
	}

	db := config.Ulbimongoconn
	adminCollection := db.Collection("Admin")
	var admin inimodel.Admin
	err = adminCollection.FindOne(context.Background(), bson.M{"email": userInfo.Email}).Decode(&admin)
	if err == mongo.ErrNoDocuments {
		// Jika admin belum ada, buat data admin baru
		admin.ID = primitive.NewObjectID()
		admin.User_name = userInfo.Name
		admin.Email = userInfo.Email
		_, err = adminCollection.InsertOne(context.Background(), admin)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Gagal menyimpan admin")
		}
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Kesalahan database")
	}

	// Menghasilkan token JWT
	jwtToken, err := iniconfig.GenerateJWT(admin)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal menghasilkan JWT")
	}

	// Mengembalikan token JWT
	return c.JSON(fiber.Map{"token": jwtToken})
}