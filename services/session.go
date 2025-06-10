package services

import (
	"database/sql"
	"encoding/base32"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db        *sql.DB
	tableName string
	Codecs    []securecookie.Codec
	Options   *sessions.Options
	keyPrefix string
}

func NewSessionStore(db *sql.DB, tableName, keyPrefix string, keyPairs ...[]byte) sessions.Store {
	return newPostgresStore(db, tableName, keyPrefix, keyPairs...)
}

// NewPostgresStore creates a new instance of PostgresStore, ensuring the table exists.
func newPostgresStore(db *sql.DB, tableName, keyPrefix string, keyPairs ...[]byte) sessions.Store {
	createTable := `
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		session_key TEXT PRIMARY KEY,
		session_data BYTEA,
		updated_at TIMESTAMP
	);`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal("Error connecting to postgres session store database")
		return nil
	}
	return &PostgresStore{
		db:        db,
		tableName: tableName,
		keyPrefix: keyPrefix,
		Codecs:    securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 1, // 1 days
		},
	}
}

// Get returns a session for the given name after adding it to the registry.
func (s *PostgresStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a new session for the given name without adding it to the registry.
func (s *PostgresStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true

	cookie, err := r.Cookie(name)
	if err != nil {
		// No cookie means new session
		return session, nil
	}

	sessionID := cookie.Value
	if sessionID == "" {
		return session, nil
	}

	var data []byte
	query := `SELECT session_data FROM ` + s.tableName + ` WHERE session_key = $1`
	err = s.db.QueryRow(query, s.keyPrefix+sessionID).Scan(&data)
	if err == sql.ErrNoRows {
		return session, nil // Not found: new session
	} else if err != nil {
		return session, err
	}

	// Decode session.Values
	if err = securecookie.DecodeMulti(name, string(data), &session.Values, s.Codecs...); err == nil {
		session.IsNew = false
	} else {
		return session, err
	}

	return session, nil
}

// Save adds a single session to the response.
func (s *PostgresStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge < 0 {
		// Delete session
		s.delete(session)
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values, s.Codecs...)
	if err != nil {
		return err
	}

	sessionID := session.ID
	if sessionID == "" {
		// Generate a new session ID
		sessionID = strings.TrimRight(base32.StdEncoding.EncodeToString(
			securecookie.GenerateRandomKey(32)), "=")
		session.ID = sessionID
	}

	_, err = s.db.Exec(
		`INSERT INTO `+s.tableName+` (session_key, session_data, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_key) DO UPDATE SET session_data=$2, updated_at=$3`,
		s.keyPrefix+sessionID, []byte(encoded), time.Now(),
	)
	if err != nil {
		return err
	}

	// Set cookie with session ID
	http.SetCookie(w, sessions.NewCookie(session.Name(), sessionID, session.Options))
	return nil
}

// Deletes the session from the database.
func (s *PostgresStore) delete(session *sessions.Session) error {
	_, err := s.db.Exec(
		`DELETE FROM `+s.tableName+` WHERE session_key = $1`,
		s.keyPrefix+session.ID,
	)
	return err
}

// // Utility: Register types that will be saved in session.Values
// func init() {
// 	gob.Register(map[interface{}]interface{}{})
// }
