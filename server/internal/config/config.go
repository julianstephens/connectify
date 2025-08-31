package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port              string `envconfig:"PORT" default:"8080"`
	DBPassword        string `envconfig:"DB_PASSWORD"`
	LogType           string `envconfig:"LOG_TYPE" default:"zap"`
	Auth0Domain       string `envconfig:"AUTH0_DOMAIN"`
	Auth0ClientID     string `envconfig:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret string `envconfig:"AUTH0_CLIENT_SECRET"`
	Auth0Audience     string `envconfig:"AUTH0_AUDIENCE"`
	LogFormat         string `envconfig:"LOG_FORMAT" default:"json"` // "json" or "human"
}

var AppConfig Config

// Load loads environment variables into AppConfig
func Load() {
	if err := envconfig.Process("", &AppConfig); err != nil {
		log.Fatalf("Failed to load env config: %v", err)
	}
}
