package main

import (
	"log"
	"phongsathorn/go_backend_gin/database"
	"phongsathorn/go_backend_gin/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()
	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}
