package controllers

import (
	// "fmt"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/models"
	"github.com/BeeCodingAI/triana-api/schemas"
	"github.com/BeeCodingAI/triana-api/services"
	"github.com/BeeCodingAI/triana-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateSessionResponse(c *gin.Context) {
	session_id := c.Param("id")

	// check if session_id exists in the database
	var existingSesssion models.Session

	err := config.DB.Preload("User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // Order messages by created_at in ascending order (for earlier messages first)
		}).
		Where("id = ?", session_id).First(&existingSesssion).Error

	if err != nil {
		c.JSON(404, gin.H{"message": "Session not found"})
		return
	}

	// get the new message from user
	var input schemas.SessionChatInput

	// bind and validate the input
	if valid, _ := utils.BindAndValidate(c, &input); !valid {
		return // The response has already been sent in the utility function
	}

	// get the message reply from LLM
	reply, err := services.GetLLMResponse(input.NewMessage, &existingSesssion)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	var LLMResponse schemas.LLMResponse
	LLMResponse, err = services.ParseJSON(reply)

	if next_action := LLMResponse.NextAction; next_action == "CONITNUE_CHAT" {
		// update the session with the new message and LLM response
		err = services.UpdateChatHistory(session_id, input.NewMessage, LLMResponse.Reply)
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
	} else if next_action == "APPOINTMENT" {
		// create queue numbe
		var queue models.Queue
		
		// parse string to uuid
		parsedUUID, err := uuid.Parse(session_id)
		if err != nil {
			c.JSON(500, gin.H{"message": "Invalid session ID"})
			return
		}
		queue.SessionID = parsedUUID

		// parse string to doctor id
		parsedUUID, err = uuid.Parse(LLMResponse.DoctorID)
		if err != nil {
			c.JSON(500, gin.H{"message": "Invalid doctor ID"})
			return
		}

		queue.DoctorID = parsedUUID

		// create queue (number is auto incremented)
		err = config.DB.Create(&queue).Error
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to create queue"})
			return
		}
	}

	// update the chat history with the new message and LLM response
	err = services.UpdateChatHistory(session_id, input.NewMessage, LLMResponse.Reply)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// send the response back to the client
	c.JSON(200, gin.H{
		"message":     "Chat history updated successfully",
		"reply":       LLMResponse.Reply,
		"session_id":  session_id,
		"next_action": LLMResponse.NextAction,
	})
}

func GetActiveSession(c *gin.Context) {
	session_id := c.Param("id")

	var session models.Session

	session, err := services.GetSessionData(session_id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Session not found"})
		return
	}

	// make sure sessoin is active (prediagnosis is not done yet)
	if session.Prediagnosis != "" {
		c.JSON(400, gin.H{"message": "Session is completed"})
		return
	}

	// send the response back to the client
	c.JSON(200, session)
}
