package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	ServerHost       string
	ServerPort       string
	KeycloakURL      string // Authorization base url
	KCClientID       string // client id oauth
	RedirectURL      string // valid redirect url
	KCClientSecret   string // optional
	Realm            string // keycloak realm
	StorageEndpoint  string
	StorageRegion    string
	StorageAccessKey string
	StorageSecretKey string
}

func Load() *Config {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8081")
	viper.SetDefault("KEYCLOAK_URL", "http://localhost:8080")
	viper.SetDefault("STORAGE_ENDPOINT", "http://localhost:9000")
	// this doesn't do anything but is needed for s3 api call compatibility
	viper.SetDefault("STORAGE_REGION", "us-west-1")
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig() // .env is optional
	viper.AutomaticEnv()

	return &Config{
		DBHost:           viper.GetString("DB_HOST"),
		DBPort:           viper.GetString("DB_PORT"),
		DBUser:           viper.GetString("DB_USER"),
		DBPassword:       viper.GetString("DB_PASSWORD"),
		DBName:           viper.GetString("DB_NAME"),
		DBSSLMode:        viper.GetString("DB_SSL_MODE"),
		KeycloakURL:      viper.GetString("KEYCLOAK_URL"),
		KCClientID:       viper.GetString("KEYCLOAK_CLIENT_ID"),
		KCClientSecret:   viper.GetString("KEYCLOAK_CLIENT_SECRET"),
		Realm:            viper.GetString("KEYCLOAK_REALM"),
		ServerHost:       viper.GetString("SERVER_HOST"),
		ServerPort:       viper.GetString("SERVER_PORT"),
		StorageEndpoint:  viper.GetString("STORAGE_ENDPOINT"),
		StorageRegion:    viper.GetString("STORAGE_REGION"),
		StorageAccessKey: viper.GetString("STORAGE_ACCESS_KEY"),
		StorageSecretKey: viper.GetString("STORAGE_SECRET_KEY"),
	}
}
