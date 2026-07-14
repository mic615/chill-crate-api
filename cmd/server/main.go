package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/mic615/chill-crate-api/internal/auth"
	"github.com/mic615/chill-crate-api/internal/config"
	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/handlers"
	"github.com/mic615/chill-crate-api/internal/routes"
	"github.com/mic615/chill-crate-api/internal/storage"
)

func main() {
	router := gin.Default()
	cfg := config.Load()
	database.Connect(cfg)
	s := storage.Connect(cfg)
	a := auth.NewAuthenticator(cfg, database.DB)
	h := handlers.NewHandler(database.DB, s)
	routes.RegisterRoutes(router, h, a.AuthMiddleware())
	// Inform the user where the server is listening
	log.Println("Running @ http://" + cfg.ServerHost + ":" + cfg.ServerPort)
	// Print out and exit(1) to the OS if the server cannot run
	log.Fatalln(router.Run(cfg.ServerHost + ":" + cfg.ServerPort))
}
