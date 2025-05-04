package services

import (
	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/models"
)

func GetAllDoctors() []models.Doctor {
	var doctors []models.Doctor
	err := config.DB.Find(&doctors).Error
	if err != nil {
		return nil
	}
	return doctors
}
