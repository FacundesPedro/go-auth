package domain

import (
	"FacundesPedro/go-auth/utils"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexedwards/scs/v2"
)

func NewSessionManager(db *sql.DB) *scs.SessionManager {

	// var SessionManager *scs.SessionManager
	// Initialize the session manager
	// this runs when package is imported
	sessionManager := scs.New()
	//
	createSessionsTable(db)
	setDefaultSessionConfig(sessionManager)
	//
	return sessionManager
}

func setDefaultSessionConfig(scs *scs.SessionManager) {
	scs.Lifetime = 24 * time.Hour
	scs.Cookie.Persist = true
	scs.Cookie.HttpOnly = true
	scs.Cookie.Secure = utils.GetEnvAsBool("HTTPS", false) // enable with HTTPS
}

func createSessionsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS sessions (
        token TEXT PRIMARY KEY,
        data BYTEA NOT NULL,
        expiry TIMESTAMPTZ NOT NULL
    );
    CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create sessions table: %v", err)
	}
	log.Println("'Sessions' table created/verified successfully")
	return nil
}
