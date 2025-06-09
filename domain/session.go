package domain

import (
	"FacundesPedro/go-auth/utils"

	"github.com/gorilla/sessions"
)

func NewSessionManager(key string) *sessions.CookieStore {

	// var SessionManager *scs.SessionManager
	// Initialize the session manager
	// this runs when package is imported
	// sessionManager := scs.New()
	cookieStore := sessions.NewCookieStore([]byte(key))
	//
	// createSessionsTable(db)
	setDefaultSessionConfig(cookieStore)
	//
	return cookieStore
}

func setDefaultSessionConfig(store *sessions.CookieStore) {
	//store.Options.MaxAge = 24 * time.Hour * time.Second
	store.Options.Path = "/"
	store.MaxAge(24 * 60 * 60) // 24 hours * 60 minutes * 60 seconds (86400 seconds)
	store.Options.HttpOnly = true
	store.Options.Secure = utils.GetEnvAsBool("HTTPS", false) // enable with HTTPS
}
