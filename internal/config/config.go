package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerHost string
	ServerPort string

	StorageEndpoint  string
	StorageRegion    string
	StorageAccessKey string
	StorageSecretKey string
}

func Load() *Config {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("SERVER_PORT", "8081")
	viper.SetDefault("STORAGE_ENDPOINT", "http://localhost:9000")
	// this doesn't do anything but is needed for s3 api call compatibility
	viper.SetDefault("STORAGE_REIGON", "us-west-1")
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig() // .env is optional
	viper.AutomaticEnv()

	return &Config{
		DBHost:           viper.GetString("DB_HOST"),
		DBPort:           viper.GetString("DB_PORT"),
		DBUser:           viper.GetString("DB_USER"),
		DBPassword:       viper.GetString("DB_PASSWORD"),
		DBName:           viper.GetString("DB_NAME"),
		DBSSLMode:        viper.GetString("DB_SSLMODE"),
		ServerHost:       viper.GetString("SERVER_HOST"),
		ServerPort:       viper.GetString("SERVER_PORT"),
		StorageEndpoint:  viper.GetString("STORAGE_ENDPOINT"),
		StorageRegion:    viper.GetString("STORAGE_REIGON"),
		StorageAccessKey: viper.GetString("STORAGE_ACCESS_KEY"),
		StorageSecretKey: viper.GetString("STORAGE_SECRET_KEY"),
	}
}
