package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	ID uint `json:"id" gorm:"primaryKey"`
}
