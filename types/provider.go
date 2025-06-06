package types

import "github.com/markbates/goth"

type ProviderConfig struct {
	Name   string
	Config goth.Provider
}
