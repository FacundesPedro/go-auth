package services

import (
	"FacundesPedro/go-auth/types"

	"github.com/markbates/goth"
)

func NewAuth(providers []types.ProviderConfig) {
	// gothic.Store = nil
	//
	goth.UseProviders(
		mapGothProviders(providers)...,
	)
	//
	// return
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
