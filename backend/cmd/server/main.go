package main

import (
	"log"

	"cocktail-bar-system/internal/config"
	"cocktail-bar-system/internal/router"
	"cocktail-bar-system/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	database.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.SetupRoutes(r)

	log.Printf("Server starting on %s:%s", config.AppConfig.ServerHost, config.AppConfig.ServerPort)
	err := r.Run(config.AppConfig.ServerHost + ":" + config.AppConfig.ServerPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
