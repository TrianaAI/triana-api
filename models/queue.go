package models

import "github.com/google/uuid"

type Queue struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DoctorID  string    `json:"doctor_id" gorm:"type:uuid;not null"`
	SessionID string    `json:"session_id" gorm:"type:uuid;not null"`
	CreatedAt string    `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt string    `json:"updated_at" gorm:"type:timestamp;not null"`
}
