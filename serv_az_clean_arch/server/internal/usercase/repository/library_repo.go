package repository

import (
	"context"
	"server/internal/domain/models"
)

type LibraryRepository interface {
	GetAllWords() ([]*models.Library, error)
	GetTranslationRus(word string) ([]*models.Library, error)
	GetTranslationRusLike(word string) ([]*models.Library, error)
	GetTranslationRusLikeWord(word string) (*models.Library, error)
	GetTranslationEngl(word string) ([]*models.Library, error)
	GetTranslationEnglLike(word string) ([]*models.Library, error)
	GetTranslationEnglLikeWord(word string) (*models.Library, error)
	InsertWordsLibrary(ctx context.Context, library []*models.Library) error
	InsertWordLibrary(ctx context.Context, word *models.Library) error
	UpdateWord(ctx context.Context, word *models.Library) error
	InitWordsMap() error
	UpdateWordsMap() error
	GetAllTopics() ([]string, error)
}
