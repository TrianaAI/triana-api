package main

import (
	"github.com/gin-gonic/gin"

	"github.com/BeeCodingAI/triana-api/config"
)

func main() {
	// Connect to the database
	config.ConnectDatabase()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(":8080") // listen and serve on
}
