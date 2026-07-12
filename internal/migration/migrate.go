package main

import (
	"github.com/mic615/chill-crate-api/internal/config"
	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/models"
)

func main() {
	database.Connect(config.Load())
	err := database.DB.AutoMigrate(
		&models.Group{},
		&models.Membership{},
		&models.Bucket{},
		&models.Object{},
		&models.User{},
	)
	if err != nil {
		panic(err)
	}
}
