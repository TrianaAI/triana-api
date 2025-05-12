package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/controllers"
)

func main() {
	// load .env file
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else if os.IsNotExist(err) {
		log.Println(".env file not found, using environment variables from Docker")
	} else {
		log.Fatal("Error checking .env file:", err)
	}

	// Connect to the database
	config.ConnectDatabase()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	r := gin.Default()
	r.Use(cors.New(corsConfig))

	// register routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/verify-otp", controllers.VerifyOTP)

	// session routes
	r.GET("/session/:id", controllers.GetActiveSession)
	r.POST("/session/:id", controllers.GenerateSessionResponse)
	r.POST("/session/:id/diagnose", controllers.DoctorDiagnose)

	// queue routes
	r.GET("/queue/:doctor_id", controllers.GetCurrentQueue)

	// test routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(":8080") // listen and serve on
}
