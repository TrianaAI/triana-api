package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/models"
	"github.com/BeeCodingAI/triana-api/schemas"
	"github.com/BeeCodingAI/triana-api/utils"
	"google.golang.org/genai"
)

func convertMessageToGenaiContent(message schemas.Message) *genai.Content {
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

func GetLLMResponse(newMessage string, session models.Session) (string, error) {
	// initialize the Gemini client with your API key and backend
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	// from the chat history, unmarshal the JSON to a slice of Message structs
	var storedHistory []schemas.Message
	err := json.Unmarshal(session.ChatHistory, &storedHistory)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling chat history: %w", err)
	}

	// build the genai history
	var genaiHistory []*genai.Content

	// build the system prompt using the session data
	systemPromptText := buildSystemPrompt(session)
	log.Printf("System Prompt: %s\n", systemPromptText)

	// add system prompt to the history
	systemPrompt := genai.NewContentFromText(systemPromptText, genai.RoleUser)
	genaiHistory = append(genaiHistory, systemPrompt)

	// add the stored messages to the history
	for _, messageItem := range storedHistory {
		content := convertMessageToGenaiContent(messageItem)
		if content != nil {
			genaiHistory = append(genaiHistory, content)
		}
	}

	chat, _ := client.Chats.Create(ctx, "gemini-2.0-flash", nil, genaiHistory)
	res, _ := chat.SendMessage(ctx, genai.Part{Text: newMessage})

	// get the response from the LLM
	if len(res.Candidates) > 0 {
		return res.Candidates[0].Content.Parts[0].Text, nil
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

	// unmarshal the chat history to a slice of Message structs
	var chatHistory []schemas.Message
	err = json.Unmarshal(session.ChatHistory, &chatHistory)
	if err != nil {
		return fmt.Errorf("error unmarshalling chat history: %v", err)
	}

	// append the new message and LLM response to the chat history
	newUserMessage := schemas.Message{Role: "user", Content: newMessage}
	newLLMResponse := schemas.Message{Role: "triana", Content: LLMResponse}
	chatHistory = append(chatHistory, newUserMessage, newLLMResponse)

	// marshal the updated chat history back to JSON
	updatedChatHistory, err := json.Marshal(chatHistory)
	if err != nil {
		return fmt.Errorf("error marshalling updated chat history: %v", err)
	}

	// update the session's chat history in the database
	session.ChatHistory = updatedChatHistory
	err = config.DB.Save(&session).Error
	if err != nil {
		return fmt.Errorf("error updating session: %v", err)
	}

	return nil
}

func buildSystemPrompt(session models.Session) string {
	// Build the system prompt using the user's data
	userDataText := fmt.Sprintf(
		"\nHere's the user's data: \n\nName:%s\nAge:%s\nNationality:%s\nWeight: %d\nHeight: %d\nHeartrate: %d\nBodytemp: %f\n",
		session.User.Name,
		utils.DateToAgeString(session.User.DOB),
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
	doctorListText := fmt.Sprintf("\nHere are the doctors available [ID] Name (Specialty):\n%s", doctorList)
	systemPromptText := fmt.Sprintf(os.Getenv("TRIANA_SYS_PROMPT"), userDataText, doctorListText)
	return systemPromptText
}
