package server

import "server/internal/domain/models"

type Rsvp struct {
	Words    []*models.Library
	WordRus  string
	WordEng  string
	Word     string
	Quantity int
}
