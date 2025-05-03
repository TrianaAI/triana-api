package controllers

import (
	"github.com/BeeCodingAI/triana-api/schemas"
	"github.com/BeeCodingAI/triana-api/services"
	"github.com/BeeCodingAI/triana-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func RegisterUser(c *gin.Context) {
	var input schemas.RegisterUserInput

	// bind and validate the request body to the input struct
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // the response has already been sent in the utility function
	}

	// Call the service to register the user
	user, err := services.RegisterUser(input)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// Send user id, name, and email in the response to prompt OTP verification
	c.JSON(200, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func VerifyOTP(c *gin.Context) {
	var input schemas.OTPInput

	// bind and validate the request body to the input struct
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// Call the service to verify the OTP
	session, err := services.ValidateOTP(input)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "OTP verified successfully", "session": session})
}
