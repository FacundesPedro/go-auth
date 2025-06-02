package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"FacundesPedro/go-auth/constants"
	"FacundesPedro/go-auth/domain"
	"FacundesPedro/go-auth/dto"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
)

type App struct {
	config         *oauth2.Config
	sessionManager *scs.SessionManager
}

// NewApp constructs the application handler
func NewApp(config *oauth2.Config, sessionManager *scs.SessionManager) *App {
	return &App{config: config, sessionManager: sessionManager}
}

// LoginHandler serves the login page
func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./public/index.html")
	if err != nil {
		log.Print(err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// OAuthHandler redirects to Google OAuth consent screen
func (a *App) OAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// OAuthCallbackHandler exchanges code, fetches user info, and stores session
func (a *App) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := a.config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange error", http.StatusBadRequest)
		return
	}

	// Store token in session for later refresh
	a.sessionManager.Put(r.Context(), "oauthToken", token)

	client := a.config.Client(context.Background(), token)
	resp, err := client.Get(constants.GOOGLE_V2_USERINFO_URL)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	var info dto.AUTH_USERINFO_V2
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		http.Error(w, "Decode user info error", http.StatusInternalServerError)
		return
	}

	user := domain.LocalUser{Name: info.Name, Email: info.Email, Image: info.Picture}
	a.sessionManager.Put(r.Context(), "user", user)

	fmt.Fprintf(w, "Hello, %s!", user.Name)
}

// RefreshHandler refreshes the OAuth token and updates the session
func (a *App) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	raw := a.sessionManager.Get(ctx, "oauthToken")
	token, ok := raw.(*oauth2.Token)
	if !ok || token == nil {
		http.Error(w, "No token in session", http.StatusUnauthorized)
		return
	}

	ts := a.config.TokenSource(context.Background(), token)
	newToken, err := ts.Token()
	if err != nil {
		http.Error(w, "Token refresh failed", http.StatusInternalServerError)
		return
	}
	a.sessionManager.Put(ctx, "oauthToken", newToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "token refreshed"})
}
