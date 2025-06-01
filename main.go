package main

import (
	"FacundesPedro/go-auth/domain"
	"FacundesPedro/go-auth/handler"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// session manager configs
	gob.Register(domain.LocalUser{})
	var sessionManager *scs.SessionManager
	//
	clientID := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	port := 5000
	//
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("http://localhost:%d/auth/callback", port),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	//
	sessionManager = scs.New()
	//
	app := handler.NewApp(oauthConfig, sessionManager)
	router := http.NewServeMux()
	//
	router.HandleFunc("GET /auth/oauth", app.OAuthHandler)
	router.HandleFunc("GET /auth/callback", app.OAuthCallbackHandler)
	router.HandleFunc("GET /auth/login", app.LoginHandler)
	//
	http.ListenAndServe(fmt.Sprintf(":%d", port), sessionManager.LoadAndSave(router))
}
