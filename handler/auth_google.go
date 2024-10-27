package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rayfanaqbil/Zenverse-BP/config"
	iniconfig "github.com/rayfanaqbil/zenverse-BE/v2/config"
	inimodel "github.com/rayfanaqbil/zenverse-BE/v2/model"
	inimodul "github.com/rayfanaqbil/zenverse-BE/v2/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// GoogleLogin memulai alur autentikasi Google OAuth
func GoogleLogin(c *fiber.Ctx) error {
	url := inimodul.GoogleOAuthConfig.AuthCodeURL("state-token")
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).SendString("Kode tidak ditemukan")
	}

	token, err := inimodul.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal menukar token")
	}

	client := inimodul.GoogleOAuthConfig.Client(context.Background(), token)

	oauth2Service, err := oauth2.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal membuat layanan OAuth2")
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Gagal mengambil info pengguna")
	}

	if !inimodul.IsAdmin(config.Ulbimongoconn, "admin", userInfo.Email) {
		return c.Status(fiber.StatusForbidden).SendString("Akses hanya untuk admin")
	}

	db := config.Ulbimongoconn
	adminCollection := db.Collection("admin")
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
