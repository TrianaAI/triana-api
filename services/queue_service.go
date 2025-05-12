package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/models"
	"github.com/BeeCodingAI/triana-api/schemas"
	"github.com/google/uuid"
)

func GenerateQueue(sessionID string, doctorID string) (*models.Queue, error) {
	var queue models.Queue

	// parse the sessionID and doctorID to UUID
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}
	queue.SessionID = sessionUUID

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID: %w", err)
	}
	queue.DoctorID = doctorUUID

	// get start of the current day
	todayStart := time.Now().Truncate(24 * time.Hour)

	// find the latest queue entry today universally
	var latestQueue models.Queue
	err = config.DB.
		Where("doctor_id = ?", doctorID).
		Where("created_at >= ?", todayStart).
		Order("number DESC").
		First(&latestQueue).Error

	if err != nil {
		// no queue entries found today, start from 1
		queue.Number = 1
	} else {
		// increment the latest queue number by 1
		queue.Number = latestQueue.Number + 1
	}

	// set the created and updated time
	now := time.Now()
	queue.CreatedAt = now
	queue.UpdatedAt = now

	// insert the queue entry into the database
	err = config.DB.Create(&queue).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create queue entry: %w", err)
	}

	return &queue, nil
}

func GetCurrentQueue(doctorID uuid.UUID) (*models.Queue, error) {
	todayStart := time.Now().Truncate(24 * time.Hour)

	var queue models.Queue
	err := config.DB.
		Where("doctor_id = ?", doctorID).
		Where("created_at >= ?", todayStart).
		Order("number ASC").
		First(&queue).Error

	if err != nil {
		return nil, fmt.Errorf("no queue found for today: %w", err)
	}

	return &queue, nil
}

func SendQueueEmail(to string, queue int, currentQueue int, token string) (map[string]interface{}, error) {
	email := schemas.Email{
		To:      to,
		Subject: "Queue Notification",
		Body:    fmt.Sprintf("Your queue number is %d. The current queue number is %d.", queue, currentQueue),
		From:    "test@example.com",
		HTML:    "test",
	}

	// URL of the email service
	url := "http://52.230.88.220:16250/send-email"

	jsonData, err := json.Marshal(email)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error from email service: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func injectQueueIntoHTML(queue int, currentQueue int, doctor models.Doctor) string {
	// Path to the HTML file
	htmlFilePath := "emails/queue_mail.html"

	// Read the HTML file
	htmlBytes, err := os.ReadFile(htmlFilePath)
	if err != nil {
		log.Fatalf("Failed to read HTML file: %v", err)
	}

	// Convert the file content to a string
	htmlString := string(htmlBytes)

	htmlString = strings.Replace(htmlString, "{{queue}}", fmt.Sprintf("%d", queue), -1)
	htmlString = strings.Replace(htmlString, "{{current_queue}}", fmt.Sprintf("%d", currentQueue), -1)
	htmlString = strings.Replace(htmlString, "{{doctor_name}}", doctor.Name, -1)
	htmlString = strings.Replace(htmlString, "{{doctor_speciality}}", doctor.Specialty, -1)
	htmlString = strings.Replace(htmlString, "{{room_number}}", doctor.Roomno, -1)

	return htmlString
}
