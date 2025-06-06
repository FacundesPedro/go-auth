package config

import (
	"FacundesPedro/go-auth/utils"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	AppBaseURL               string
	Port                     string
	GoogleClientID           string
	GoogleSecretID           string
	GoogleCallbackURL        string
	PostgresStringConnection string
}

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		GoogleClientID:           utils.GetEnvOrError("GOOGLE_OAUTH_CLIENT_ID"),
		GoogleSecretID:           utils.GetEnvOrError("GOOGLE_OAUTH_CLIENT_SECRET"),
		GoogleCallbackURL:        utils.GetEnvOrError("GOOGLE_OAUTH_CALLBACK_URL"),
		AppBaseURL:               utils.GetEnv("APP_BASE_URL", "http://localhost"),
		PostgresStringConnection: utils.GetEnvOrError("POSTGRES_STRING_CONNECTION"),
		Port:                     utils.GetEnv("APP_PORT", "5000"),
	}
}
