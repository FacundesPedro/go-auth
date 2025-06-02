package main

import (
	"FacundesPedro/go-auth/domain"
	"FacundesPedro/go-auth/handler"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// session manager configs
	// var sessionManager *scs.SessionManager
	gob.Register(domain.LocalUser{})
	gob.Register(&oauth2.Token{})
	//
	baseURL := func() string {
		if v := os.Getenv("APP_BASE_URL"); v != "" {
			return v
		}
		return "http://localhost"
	}()
	clientID := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	port := 5000
	//
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  fmt.Sprintf("%s:%d/auth/callback", baseURL, port),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	// Initialize the session manager
	sessionManager := scs.New()
	sessionManager.Lifetime = 5 * time.Second
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = true // Enable when using HTTPS
	//
	// sessionManager.Store = redisstore.New(redisClient)
	//
	app := handler.NewApp(oauthConfig, sessionManager)
	router := http.NewServeMux()
	//
	router.HandleFunc("GET /auth/oauth", app.OAuthHandler)
	router.HandleFunc("GET /auth/callback", app.OAuthCallbackHandler)
	router.HandleFunc("GET /auth/login", app.LoginHandler)
	// router.HandleFunc("/auth/refresh", app.RefreshHandler)
	//
	log.Printf("Server starting on %s:%d", baseURL, port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), sessionManager.LoadAndSave(router))
}
