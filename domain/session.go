package domain

import (
	"FacundesPedro/go-auth/utils"

	"github.com/gorilla/sessions"
)

func SetDefaultSessionConfig(store *sessions.CookieStore) {
	//store.Options.MaxAge = 24 * time.Hour * time.Second
	store.Options.Path = "/"
	store.MaxAge(24 * 60 * 60) // 24 hours * 60 minutes * 60 seconds (86400 seconds)
	store.Options.HttpOnly = true
	store.Options.Secure = utils.GetEnvAsBool("HTTPS", false) // enable with HTTPS
}
