package main

import (
	"github.com/gin-gonic/gin"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/controllers"
)

func main() {
	// Connect to the database
	config.ConnectDatabase()

	r := gin.Default()

	// register routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/verify-otp", controllers.VerifyOTP)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(":8080") // listen and serve on
}
