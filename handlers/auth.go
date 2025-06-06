package handlers

import (
	"FacundesPedro/go-auth/constants"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

// LoginHandler serves the login page
func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./public/index.html")
	if err != nil {
		log.Print(err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (a *App) ProviderHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = gothic.GetContextWithProvider(r, provider)
	//
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		message := constants.USER_ALREADY_AUTHENTICATED(user.Name)
		//
		log.Print(message)
		fmt.Fprint(w, message)
		//
		return
	}
	//
	gothic.BeginAuthHandler(w, r)
}

func (a *App) HandleProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = gothic.GetContextWithProvider(r, provider)

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	//
	// a.sessionManager.Put(r.Context(), "user", user)
	log.Print(user.Name)
	log.Print(constants.USER_SUCCESS_AUTHENTICATED(user.Name))
	//
	//w.WriteHeader(http.StatusTemporaryRedirect)
	//w.Header().Set("Location", "/")
}
