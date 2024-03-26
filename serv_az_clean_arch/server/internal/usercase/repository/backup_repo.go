package repository

import (
	"os"
	"server/internal/domain/models"
)

type BackUpCopyRepo interface {
	GetAllWordsFromBackUpXlsx() ([]*models.Library, error)
	SaveWordsAsXLSX(words []*models.Library) error
	OpenFile() (*os.File, error)
}
