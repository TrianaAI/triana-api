package models

type Doctor struct {
	ID    string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name  string `json:"name" gorm:"type:varchar(100);not null"`
	Email string `json:"email" gorm:"type:varchar(100);unique;not null"`
}
