package handlers

import (
	"FacundesPedro/go-auth/types"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
)

type App struct {
	sessionManager *scs.SessionManager
	providers      []types.ProviderConfig
}

type AppGoogle struct {
	config         *oauth2.Config
	sessionManager *scs.SessionManager
}

// NewApp constructs the application handler
func NewApp(providers []types.ProviderConfig, sessionManager *scs.SessionManager) *App {
	return &App{providers: providers, sessionManager: sessionManager}
}

// NewApp constructs the application handler
func NewAppGoogle(config *oauth2.Config, sessionManager *scs.SessionManager) *AppGoogle {
	return &AppGoogle{config: config, sessionManager: sessionManager}
}
