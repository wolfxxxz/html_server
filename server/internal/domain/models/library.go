package models

import (
	"gorm.io/gorm"
)

type Library struct {
	gorm.Model
	ID int `json:"ID" gorm:"primaryKey"`
	//ID            int       `json:"id" `
	English       string    `json:"english"`
	Russian       string    `json:"russian"`
	Theme         string    `json:"theme"`
	PartsOfSpeech string    `json:"part_of_speech"`
	Phrases       []*Phrase `gorm:"many2many:library_phrases;" json:"library_phrases"`
	Exceptions    string    `json:"exceptions"`
}

type Phrase struct {
	gorm.Model
	ID        int       `json:"id" gorm:"primaryKey"`
	English   string    `json:"english"`
	Russian   string    `json:"russian"`
	Libraries []Library `gorm:"many2many:library_phrases;" json:"libraries"`
}
