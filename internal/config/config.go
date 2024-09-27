package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var google_clientid = os.Getenv("GOOGLE_CLIENT_ID")
var google_clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
var Google_conf = &oauth2.Config{
	ClientID:     google_clientid,
	ClientSecret: google_clientSecret,
	RedirectURL:  "http://localhost:8080/api/v1/auth/callback/google",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}
var Frontend_home = os.Getenv("FRONTEND_HOME")
var Frontend_password = os.Getenv("FRONTEND_PASSWORD")
