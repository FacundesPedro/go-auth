package services

import (
	"FacundesPedro/go-auth/domain"
	"FacundesPedro/go-auth/types"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

const (
	SessionField string = "_session"
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
func (s *AuthService) GetSessionUser(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, SessionField)
	if err != nil {
		return goth.User{}, err
	}
	//
	u := session.Values["user"]
	//
	if u == nil {
		return goth.User{}, fmt.Errorf("user is not authenticated! %v", u)
	}

	return u.(goth.User), nil
}
func (s *AuthService) RemoveUserSession(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, SessionField)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = goth.User{}
	// delete the cookie immediately
	session.Options.MaxAge = -1

	session.Save(r, w)
}
func (s *AuthService) RequireAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		log.Printf("user is authenticated! user: %v!", session.FirstName)

		handlerFunc(w, r)
	}
}

// utils
func mapGothProviders(providers []types.ProviderConfig) []goth.Provider {
	gothProviders := make([]goth.Provider, len(providers))
	//
	for i, p := range providers {
		gothProviders[i] = p.Config
	}
	//
	return gothProviders
}
