package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mic615/chill-crate-api/internal/config"
	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/routes"
)

func main() {
	router := gin.Default()
	cfg := config.Load()
	database.Connect(cfg)
	routes.RegisterRoutes(router)
	// Inform the user where the server is listening
	log.Println("Running @ http://" + cfg.ServerHost + ":" + cfg.ServerPort)
	// Print out and exit(1) to the OS if the server cannot run
	log.Fatalln(router.Run(cfg.ServerHost + ":" + cfg.ServerPort))

}
