package domain

import (
	"time"

	"github.com/alexedwards/scs/v2"
)

func NewSessionManager() *scs.SessionManager {

	// var SessionManager *scs.SessionManager
	// Initialize the session manager
	// this runs when package is imported
	SessionManager := scs.New()
	setDefaultSessionConfig(SessionManager)

	// if you want Redis backing:
	// redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	// Manager.Store = redisstore.New(redisClient)
	return SessionManager
}

func setDefaultSessionConfig(scs *scs.SessionManager) {
	scs.Lifetime = 24 * time.Hour
	scs.Cookie.Persist = true
	scs.Cookie.Secure = true // enable with HTTPS
}
