package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key    string = "ee5d7a03a32f3b35cbf1349eadce7049a53d98c53d9ef2828f94b1d58a7d02dd"
	MaxAge        = 86400 * 30
	IsProd        = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_ID")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	port := os.Getenv("PORT")
	// lista de provedores que queremos colocar OAUTH na nossa aplicação
	// https://github.com/markbates/goth/blob/master/examples/main.go
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, fmt.Sprintf("http://localhost:%d/auth/google/callback", port)),
	)
}
