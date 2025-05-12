package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/models"
	"github.com/BeeCodingAI/triana-api/schemas"
	"github.com/BeeCodingAI/triana-api/utils"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

func convertMessageToGenaiContent(message models.Message) *genai.Content {
	// Determine the role of the message
	var role genai.Role
	if message.Role == "user" {
		role = genai.RoleUser
	} else if message.Role == "triana" {
		role = genai.RoleModel
	} else {
		return nil // Invalid role, return nil or handle error as needed
	}

	// Create a new genai.Content object from the message content and role
	content := genai.NewContentFromText(message.Content, role)
	return content
}

func GetLLMResponse(newMessage string, session *models.Session) (string, error) {
	// initialize the Gemini client with your API key and backend
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		log.Printf("Error creating Gemini client: %v\n", err)
		return "", fmt.Errorf("error creating LLM client: %w", err)
	}

	// the chat history is stored as a one-to-many relationship in the database
	var storedHistory []models.Message = session.Messages

	// build the genai history
	var genaiHistory []*genai.Content

	// build the system prompt using the session data
	systemPromptText := buildSystemPrompt(session)
	log.Printf("System Prompt: %s\n", systemPromptText)

	var temperature float32 = 0.8
	var TopP float32 = 0.95
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemPromptText, genai.RoleUser),
		ResponseMIMEType:  "application/json",
		TopP:              &TopP,
		Temperature:       &temperature,
		MaxOutputTokens:   8192,
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"next_action":  {Type: genai.TypeString, Enum: []string{"CONTINUE_CHAT", "APPOINTMENT"}},
				"reply":        {Type: genai.TypeString},
				"doctor_id":    {Type: genai.TypeString},
				"prediagnosis": {Type: genai.TypeString},
			},
			Required: []string{"next_action", "reply", "doctor_id", "prediagnosis"},
		},
	}

	// add the stored messages to the history
	for _, messageItem := range storedHistory {
		content := convertMessageToGenaiContent(messageItem)
		if content != nil {
			genaiHistory = append(genaiHistory, content)
		}
	}

	chat, err := client.Chats.Create(ctx, os.Getenv("GEMINI_MODEL"), config, genaiHistory)
	if err != nil {
		log.Printf("Error creating chat: %v\n", err)
		return "", fmt.Errorf("error creating chat: %w", err)
	}

	res, err := chat.SendMessage(ctx, genai.Part{Text: newMessage})
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
		return "", fmt.Errorf("error sending message: %w", err)
	}

	// get the response from the LLM
	if res != nil && len(res.Candidates) > 0 && res.Candidates[0].Content != nil &&
		len(res.Candidates[0].Content.Parts) > 0 {
		text := res.Candidates[0].Content.Parts[0].Text
		return text, nil
	}

	return "", fmt.Errorf("no response from LLM")
}

func UpdateChatHistory(sessionId string, newMessage string, LLMResponse string) error {

	// get the session from the database
	var session models.Session
	err := config.DB.Where("id = ?", sessionId).First(&session).Error
	if err != nil {
		return fmt.Errorf("error fetching session: %v", err)
	}

	// append the new message and LLM response to the chat history
	now := time.Now()
	newUserMessage := models.Message{Role: "user", Content: newMessage, SessionID: session.ID, CreatedAt: now, UpdatedAt: now}
	newLLMResponse := models.Message{Role: "triana", Content: LLMResponse, SessionID: session.ID, CreatedAt: now.Add(time.Millisecond), UpdatedAt: now.Add(time.Millisecond)}

	// save the new messages to the database
	if err := config.DB.Create(&newUserMessage).Error; err != nil {
		return fmt.Errorf("error saving user message: %v", err)
	}

	if err := config.DB.Create(&newLLMResponse).Error; err != nil {
		return fmt.Errorf("error saving LLM response: %v", err)
	}

	return nil
}

func buildSystemPrompt(session *models.Session) string {
	// Build the system prompt using the user's data
	userDataText := fmt.Sprintf(
		"\nHere's the user's data: \n\nName:%s\nAge:%s\nGender:%s\nNationality:%s\nWeight: %d\nHeight: %d\nHeartrate: %d\nBodytemp: %f\n",
		session.User.Name,
		utils.DateToAgeString(session.User.DOB),
		session.User.Gender,
		session.User.Nationality,
		session.Weight,
		session.Height,
		session.Heartrate,
		session.Bodytemp,
	)

	doctors := GetAllDoctors()

	// Convert the doctors to a string representation
	var doctorList []string
	for _, doctor := range doctors {
		doctorList = append(doctorList, fmt.Sprintf("- [%s] %s (%s)\n", doctor.ID, doctor.Name, doctor.Specialty))
	}
	doctorListText := fmt.Sprintf("\nHere are the doctors available [ID] Name (Specialty):\n%s", strings.Join(doctorList, ""))

	// Get history of sessions
	var history []models.Session = GetHistory(session)

	// Convert history of sessions to a string representation
	var historyList []string
	if len(history) > 0 {
		for _, sessionItem := range history {
			historyList = append(historyList, fmt.Sprintf(
				"[%s]\nWeight: %f\nHeight: %f\nHeartrate: %f\nBodytemp: %f\nPrediagnosis: %s\n",
				sessionItem.CreatedAt.Format("2006-01-02 15:04:05"),
				sessionItem.Weight,
				sessionItem.Height,
				sessionItem.Heartrate,
				sessionItem.Bodytemp,
				sessionItem.Prediagnosis,
			))
		}
	} else {
		historyList = append(historyList, "No previous sessions found.\n")
	}
	historyListText := fmt.Sprintf("\nHere are the previous sessions:\n%s", strings.Join(historyList, ""))

	// Build the system prompt text
	systemPromptText := fmt.Sprintf("%s %s %s %s\nCurrent Time: %s", os.Getenv("TRIANA_SYS_PROMPT"), userDataText, doctorListText, historyListText, time.Now().Format("2006-01-02 15:04:05"))

	return systemPromptText
}

func GetSessionData(sessionId string) (models.Session, error) {
	// check if session_id exists in the database
	var session models.Session

	err := config.DB.
		Preload("User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // Order messages by created_at in descending order (for latest messages first)
		}).
		Where("id = ?", sessionId).First(&session).Error

	if err != nil {
		return models.Session{}, fmt.Errorf("session not found: %w", err)
	}

	return session, nil
}

func GetHistory(session *models.Session) []models.Session {
	var history []models.Session
	err := config.DB.Where("user_id = ?", session.User.ID).Where("id != ?", session.ID).Order("created_at DESC").Find(&history).Error
	if err != nil {
		log.Printf("Error fetching session history: %v\n", err)
		return []models.Session{} // Return an empty slice if there's an error
	}

	return history
}

func DoctorDiagnose(sessionId string, diagnosis string) error {
	// Fetch the session from the database
	var session models.Session
	err := config.DB.Where("id = ?", sessionId).First(&session).Error
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Update the session with the doctor's diagnosis
	session.DoctorDiagnosis = diagnosis
	session.UpdatedAt = time.Now()

	if err := config.DB.Save(&session).Error; err != nil {
		return fmt.Errorf("failed to save diagnosis: %w", err)
	}

	return nil
}

// RemoveMarkdownAndExtractJSON removes Markdown syntax and extracts the JSON content
func ParseJSON(input string) (schemas.LLMResponse, error) {

	// Extract JSON content
	var responseJSON schemas.LLMResponse
	err := json.Unmarshal([]byte(input), &responseJSON)
	if err != nil {
		return schemas.LLMResponse{}, fmt.Errorf("failed to extract JSON: %w", err)
	}

	return responseJSON, nil
}
