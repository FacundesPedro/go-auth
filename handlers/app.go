package handlers

import (
	"FacundesPedro/go-auth/services"
	"FacundesPedro/go-auth/types"
)

type App struct {
	auth      *services.AuthService
	providers []types.ProviderConfig
}

// NewApp constructs the application handler
func NewApp(providers []types.ProviderConfig, auth *services.AuthService) *App {
	return &App{providers: providers, auth: auth}
}
