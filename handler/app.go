package handler

import (
	"FacundesPedro/go-auth/constants"
	"FacundesPedro/go-auth/domain"
	"FacundesPedro/go-auth/dto"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
)

type App struct {
	config         *oauth2.Config
	client         *http.Client
	sessionManager *scs.SessionManager
}

func NewApp(config *oauth2.Config, sessionManager *scs.SessionManager) *App {
	return &App{
		config:         config,
		sessionManager: sessionManager,
	}
}

// OAuth handler create an unique OAuth URL using the client ID and
// client secret.
// then it redirect the user to the OAuth provider website to
// complete the login.
func (a *App) OAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := a.config.AuthCodeURL("hello world", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// OAuth callback handler handle the redirect request from the OAuth provider.
// It read the code query parameter and exchange it to get the access token.
// Then this handler call user info endpoint to get the user public detail eg.,
// Name, Email, Profile picture etc.
func (a *App) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	//
	t, err := a.config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//
	a.client = a.config.Client(context.Background(), t)
	//
	resp, err := a.client.Get(constants.GOOGLE_V2_USERINFO_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//
	defer resp.Body.Close()
	var v2_userinfo dto.AUTH_USERINFO_V2
	//
	if err = json.NewDecoder(resp.Body).Decode(&v2_userinfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// no error, just use the v2_userinfo
	localUser := domain.LocalUser{
		Name:  v2_userinfo.Name,
		Email: v2_userinfo.Email,
		Image: v2_userinfo.Picture,
	}
	// set cookies
	a.sessionManager.Put(r.Context(), "user", localUser)
	// return to response
	fmt.Fprintf(w, "Nome: %s \nImage: %s", localUser.Name, localUser.Image)
}

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./public/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	t.Execute(w, nil)
}
