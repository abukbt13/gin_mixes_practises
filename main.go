package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"practise/config"
	"practise/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDatabase()

	// Create a new default Gin router
	r := gin.Default()

	// Map URL routes using a custom function router.mapUrl
	router.MapUrl(r)

	// Start the server and listen for incoming HTTP requests
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
