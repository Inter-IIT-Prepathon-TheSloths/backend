package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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

var github_clientid = os.Getenv("GITHUB_CLIENT_ID")
var github_clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
var Github_conf = &oauth2.Config{
	ClientID:     github_clientid,
	ClientSecret: github_clientSecret,
	RedirectURL:  "http://localhost:8080/api/v1/auth/callback/github",
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
}
