package services

import (
	"fmt"
	"time"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/models"
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
