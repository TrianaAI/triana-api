package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/BeeCodingAI/triana-api/config"
	"github.com/BeeCodingAI/triana-api/models"
	"github.com/BeeCodingAI/triana-api/schemas"
)

func generateOTP() string {
	// Generate a random 6-digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}

func ValidateOTP(input schemas.OTPInput) (*models.Session, error) {
	// get the user from the input
	var user models.User
	err := config.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// check if OTP sent to the user is valid
	if user.OTP != input.OTP {
		return nil, fmt.Errorf("invalid OTP")
	}

	// create a new session for the user with the data from the input
	newSession := models.Session{
		UserID:    user.ID,
		Weight:    input.Weight,
		Height:    input.Height,
		Heartrate: input.Heartrate,
		Bodytemp:  input.Bodytemp,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// save the session to the database
	err = config.DB.Create(&newSession).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// update user's OTP to nil after successful validation
	user.OTP = ""
	err = config.DB.Save(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user OTP: %w", err)
	}

	return &newSession, nil
}
