package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       *uuid.UUID `json:"id" gorm:"primaryKey"`
	Email    string     `json:"user_email"`
	Name     string     `json:"first_name"`
	LastName string     `json:"last_name"`
	Password string     `json:"password"`
	Role     string     `json:"role"`
	Words    []*Word    `gorm:"many2many:user_words;" json:"user_words"`
	Learn    []*Word    `gorm:"many2many:user_learn;" json:"user_learn"`
	Learned  []*Word    `gorm:"many2many:user_learned;" json:"user_learned"`
}

type Word struct {
	gorm.Model
	ID            *uuid.UUID `json:"id" gorm:"primaryKey"`
	English       string     `json:"english"`
	Russian       string     `json:"russian"`
	Theme         string     `json:"theme"`
	PartsOfSpeech string     `json:"part_of_speech"`
}
