package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/controllers"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	config.ConnectDatabase()

	r := gin.Default()

	// register routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/verify-otp", controllers.VerifyOTP)

	// session routes
	r.POST("/session/:id", controllers.GenerateSessionResponse)

	// test routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(":8080") // listen and serve on
}
