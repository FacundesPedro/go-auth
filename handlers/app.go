package handlers

import (
	"FacundesPedro/go-auth/types"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type App struct {
	sessionManager *sessions.CookieStore
	providers      []types.ProviderConfig
}

type AppGoogle struct {
	config         *oauth2.Config
	sessionManager *sessions.CookieStore
}

// NewApp constructs the application handler
func NewApp(providers []types.ProviderConfig, sessionManager *sessions.CookieStore) *App {
	return &App{providers: providers, sessionManager: sessionManager}
}

// NewApp constructs the application handler
func NewAppGoogle(config *oauth2.Config, sessionManager *sessions.CookieStore) *AppGoogle {
	return &AppGoogle{config: config, sessionManager: sessionManager}
}
