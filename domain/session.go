package domain

import (
	"FacundesPedro/go-auth/utils"

	"github.com/gorilla/sessions"
)

// TODO: rewrite this to the the new postgreStore on maiN!
func SetDefaultSessionConfig(store *sessions.CookieStore) {
	//store.Options.MaxAge = 24 * time.Hour * time.Second
	store.Options.Path = "/"
	store.MaxAge(24 * 60 * 60) // 24 hours * 60 minutes * 60 seconds (86400 seconds)
	store.Options.HttpOnly = true
	store.Options.Secure = utils.GetEnvAsBool("APP_IS_HTTPS", false) // enable with HTTPS
}
