package main

import (
	"FacundesPedro/go-auth/config"
	"FacundesPedro/go-auth/domain"
	"FacundesPedro/go-auth/handlers"
	"FacundesPedro/go-auth/services"
	"FacundesPedro/go-auth/types"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	// _ "github.com/joho/godotenv/autoload"

	"github.com/markbates/goth/providers/google"
)

func main() {
	// main components
	var sessionStore sessions.Store
	var providers []types.ProviderConfig
	var env *config.Config
	var app *handlers.App
	var router *http.ServeMux
	var authService *services.AuthService
	//
	gob.Register(domain.LocalUser{})
	//
	env = config.InitConfig()
	// database for persistent things
	db := domain.InitDb("postgres", env.PostgresStringConnection)
	domain.PingDb("postgres", db)
	//
	// psqlStore := services.NewPostgresStore(db, "_session")
	// gothic.Store = psqlStore
	//
	providers = append(providers,
		types.ProviderConfig{
			Name: "google",
			Config: google.New(env.GoogleClientID, env.GoogleSecretID, env.GoogleCallbackURL,
				"profile", "email"),
		},
	)
	// Initialize the session store (can use de default one in case of test)
	sessionStore = services.NewSessionStore(db, "_session", "sess_", []byte(env.SessionKey))
	// activate goth and session store
	authService = services.NewAuth(providers, sessionStore)
	//
	app = handlers.NewApp(providers, authService)
	router = http.NewServeMux()
	//
	router.HandleFunc("GET /auth/{provider}", app.ProviderHandler)
	router.HandleFunc("GET /auth/{provider}/callback", app.HandleProviderCallback)
	router.HandleFunc("GET /login", app.LoginHandler)
	// router.HandleFunc("/auth/refresh", app.RefreshHandler)
	// defer's
	defer db.Close()
	//
	log.Printf("Server starting on %s:%s", env.AppBaseURL, env.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", env.Port), router)
}
