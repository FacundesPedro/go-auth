package handlers

import (
	"FacundesPedro/go-auth/services"
	"FacundesPedro/go-auth/types"
	"log"
)

type App struct {
	auth      *services.AuthService
	providers []types.ProviderConfig
}

// NewApp constructs the application handler
func NewApp(providers []types.ProviderConfig, auth *services.AuthService) *App {
	if len(providers) == 0 {
		log.Fatal("No providers were given")
		return nil
	}
	return &App{providers: providers, auth: auth}
}
