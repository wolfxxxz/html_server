package repository

import (
	"context"
	"server/internal/domain/models"
)

type WordsRepository interface {
	GetAllWords() ([]*models.Word, error)
	InsertWords(ctx context.Context, words []*models.Word) error
	InsertWord(ctx context.Context, word *models.Word) error
	UpdateWord(ctx context.Context, word *models.Word) error
	UpdateWords(ctx context.Context, words []*models.Word) error
	GetAllTopics() ([]string, error)
}
