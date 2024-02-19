package repositories

import (
	"context"
	"server/internal/apperrors"
	"server/internal/domain/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepoLibrary interface {
	GetAllWords() ([]*models.Library, error)
	GetTranslationRus(word string) ([]*models.Library, error)
	GetTranslationRusLike(word string) ([]*models.Library, error)
	GetTranslationEngl(word string) ([]*models.Library, error)
	GetTranslationEnglLike(word string) ([]*models.Library, error)
	InsertWordsLibrary(ctx context.Context, library []*models.Library) error
}

type repoLibrary struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewRepoLibrary(db *gorm.DB, log *logrus.Logger) RepoLibrary {
	return &repoLibrary{db: db, log: log}
}

func (rt *repoLibrary) GetAllWords() ([]*models.Library, error) {
	var words []*models.Library
	err := rt.db.Order("theme").Find(&words).Error
	if err != nil {
		appErr := apperrors.GetAllWordsLibErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

func (rt *repoLibrary) GetTranslationRus(word string) ([]*models.Library, error) {
	var words []*models.Library
	err := rt.db.Where("russian = ?", word).Find(&words).Error
	if err != nil {
		appErr := apperrors.GetTranslationRusErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

func (rt *repoLibrary) GetTranslationRusLike(word string) ([]*models.Library, error) {
	var words []*models.Library
	err := rt.db.Where("russian LIKE ?", "%"+word+"%").Find(&words).Error
	if err != nil {
		appErr := apperrors.GetTranslationRusLikeErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

func (rt *repoLibrary) GetTranslationEngl(word string) ([]*models.Library, error) {
	var words []*models.Library
	err := rt.db.Where("english = ?", word).Find(&words).Error
	if err != nil {
		appErr := apperrors.GetTranslationEnglErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

func (rt *repoLibrary) GetTranslationEnglLike(word string) ([]*models.Library, error) {
	var words []*models.Library
	err := rt.db.Where("english LIKE ?", "%"+word+"%").Find(&words).Error
	if err != nil {
		appErr := apperrors.GetTranslationEnglLikeErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

func (rt *repoLibrary) InsertWordsLibrary(ctx context.Context, library []*models.Library) error {
	for _, word := range library {
		if word == nil {
			appErr := apperrors.InsertWordsLibraryErr.AppendMessage("lib == nil")
			rt.log.Error(appErr)
			return appErr
		}

		tx := rt.db.WithContext(ctx)
		if tx.Error != nil {
			appErr := apperrors.InsertWordsLibraryErr.AppendMessage(tx.Error)
			rt.log.Error(appErr)
			return appErr
		}

		result := tx.Create(word)
		if result.Error != nil {
			appErr := apperrors.InsertWordsLibraryErr.AppendMessage(result.Error)
			rt.log.Error(appErr)
			return appErr
		}

		if result.RowsAffected == 0 {
			appErr := apperrors.InsertWordsLibraryErr.AppendMessage("no rows affected")
			rt.log.Error(appErr)
			return appErr
		}

		createdLib := &models.Library{}
		if err := tx.First(createdLib, "id = ?", word.ID).Error; err != nil {
			appErr := apperrors.InsertWordsLibraryErr.AppendMessage(err)
			rt.log.Error(appErr)
			return appErr
		}
	}

	return nil
}
