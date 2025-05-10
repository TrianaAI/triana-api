package controllers

import (
	"github.com/BeeCodingAI/triana-api/services"
	"github.com/BeeCodingAI/triana-api/utils"
	"github.com/gin-gonic/gin"
)

func DoctorDiagnose(c *gin.Context) {
	sessionId := c.Param("id")

	// Parse the diagnosis from the request body
	var input struct {
		Diagnosis string `json:"diagnosis" validate:"required"`
	}

	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// Call the service to save the diagnosis
	if err := services.DoctorDiagnose(sessionId, input.Diagnosis); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Diagnosis saved successfully"})
}
