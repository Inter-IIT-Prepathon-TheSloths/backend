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
var GoogleConf = &oauth2.Config{
	ClientID:     google_clientid,
	ClientSecret: google_clientSecret,
	RedirectURL:  "http://localhost:8080/api/v1/auth/callback/google",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

var github_clientid = os.Getenv("GITHUB_CLIENT_ID")
var github_clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
var GithubConf = &oauth2.Config{
	ClientID:     github_clientid,
	ClientSecret: github_clientSecret,
	RedirectURL:  "http://localhost:8080/api/v1/auth/callback/github",
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
}

var Port = os.Getenv("PORT")
var AppEnv = os.Getenv("APP_ENV")

var DbUrl = os.Getenv("DB_URL")
var DbName = os.Getenv("DB_NAME")
var JwtSecret = os.Getenv("JWT_SECRET")

var SmtpAppPass = os.Getenv("SMTP_APP_PASS")
var FromEmail = os.Getenv("FROM_EMAIL")
var FrontendUrl = os.Getenv("FRONTEND_URL")
var BackendUrl = os.Getenv("BACKEND_URL")
