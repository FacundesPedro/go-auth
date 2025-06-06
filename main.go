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

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func main() {
	// session manager configs
	var sessionManager *scs.SessionManager
	var providers []types.ProviderConfig
	// gob.Register(domain.LocalUser{})
	gob.Register(goth.User{})
	//
	env := config.InitConfig()
	// database for persistent things
	db := domain.InitDb("postgres", env.PostgresStringConnection)
	domain.PingDb("postgres", db)
	//
	providers = append(providers,
		types.ProviderConfig{
			Name:   "google",
			Config: google.New(env.GoogleClientID, env.GoogleSecretID, env.GoogleCallbackURL),
		},
	)
	// Initialize the session manager
	sessionManager = domain.NewSessionManager(db)
	sessionManager.Store = postgresstore.New(db)
	//
	app := handlers.NewApp(providers, sessionManager)
	router := http.NewServeMux()
	// TODO
	auth := services.NewAuth(,providers)
	//
	router.HandleFunc("GET /auth/{provider}", app.ProviderHandler)
	router.HandleFunc("GET /auth/{provider}/callback", app.HandleProviderCallback)
	router.HandleFunc("GET /auth/login", app.LoginHandler)
	// router.HandleFunc("/auth/refresh", app.RefreshHandler)
	// defer's
	defer db.Close()
	//
	log.Printf("Server starting on %s:%s", env.AppBaseURL, env.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", env.Port), sessionManager.LoadAndSave(router))
}
