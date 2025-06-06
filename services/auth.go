package services

import (
	"FacundesPedro/go-auth/types"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

type AuthService struct{}

func NewAuth(providers []types.ProviderConfig, store sessions.Store) *AuthService {
	gothic.Store = store
	//
	goth.UseProviders(
		mapGothProviders(providers)...,
	)
	//
	return &AuthService{}
}

func mapGothProviders(providers []types.ProviderConfig) []goth.Provider {
	gothProviders := make([]goth.Provider, len(providers))
	//
	for i, p := range providers {
		gothProviders[i] = p.Config
	}
	//
	return gothProviders
}
