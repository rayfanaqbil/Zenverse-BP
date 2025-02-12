package config

import(

 "github.com/gofiber/fiber/v2"
 "os"
 "golang.org/x/oauth2"
 "golang.org/x/oauth2/google"
)

var Iteung = fiber.Config{
	Prefork:       false,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "Iteung",
	AppName:       "Message Router",
}

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
)

var GHAccessToken string = os.Getenv("GH_ACCESS_TOKEN")