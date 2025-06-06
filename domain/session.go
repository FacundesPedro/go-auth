package domain

import (
	"FacundesPedro/go-auth/utils"
	"database/sql"

	"github.com/gorilla/sessions"
)

func NewSessionManager(db *sql.DB) *sessions.CookieStore {

	// var SessionManager *scs.SessionManager
	// Initialize the session manager
	// this runs when package is imported
	// sessionManager := scs.New()
	cookieStore := sessions.NewCookieStore()
	//
	// createSessionsTable(db)
	setDefaultSessionConfig(cookieStore)
	//
	return cookieStore
}

func setDefaultSessionConfig(store *sessions.CookieStore) {
	//store.Options.MaxAge = 24 * time.Hour * time.Second
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = utils.GetEnvAsBool("HTTPS", false) // enable with HTTPS
}
