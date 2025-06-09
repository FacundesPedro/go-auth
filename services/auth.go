package services

import (
	"FacundesPedro/go-auth/domain"
	"FacundesPedro/go-auth/types"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

const (
	SessionField string = "session"
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

func (s *AuthService) SaveUserSession(w http.ResponseWriter, r *http.Request, u domain.LocalUser) error {
	session, _ := gothic.Store.Get(r, SessionField)
	//
	session.Values["user"] = u
	//
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil

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
