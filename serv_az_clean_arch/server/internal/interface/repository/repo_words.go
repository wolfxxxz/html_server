package repository

import (
	"context"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/usercase/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type wordsRepository struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewWordsRepository(db *gorm.DB, log *logrus.Logger) repository.WordsRepository {
	return &wordsRepository{db: db, log: log}
}

func (rt *wordsRepository) GetAllWords() ([]*models.Word, error) {
	var words []*models.Word
	err := rt.db.Order("theme").Find(&words).Error
	if err != nil {
		appErr := apperrors.GetAllWordsLibErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

func (rt *wordsRepository) InsertWords(ctx context.Context, words []*models.Word) error {
	for _, word := range words {
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

		createdLib := &models.Word{}
		if err := tx.First(createdLib, "id = ?", word.ID).Error; err != nil {
			appErr := apperrors.InsertWordsLibraryErr.AppendMessage(err)
			rt.log.Error(appErr)
			return appErr
		}
	}

	return nil
}

func (rt *wordsRepository) InsertWord(ctx context.Context, word *models.Word) error {
	if word == nil {
		appErr := apperrors.InsertWordLibraryErr.AppendMessage("lib == nil")
		rt.log.Error(appErr)
		return appErr
	}

	tx := rt.db.WithContext(ctx)
	if tx.Error != nil {
		appErr := apperrors.InsertWordLibraryErr.AppendMessage(tx.Error)
		rt.log.Error(appErr)
		return appErr
	}

	result := tx.Create(word)
	if result.Error != nil {
		appErr := apperrors.InsertWordLibraryErr.AppendMessage(result.Error)
		rt.log.Error(appErr)
		return appErr
	}

	if result.RowsAffected == 0 {
		appErr := apperrors.InsertWordLibraryErr.AppendMessage("no rows affected")
		rt.log.Error(appErr)
		return appErr
	}

	createdWord := &models.Word{}
	if err := tx.First(createdWord, "id = ?", word.ID).Error; err != nil {
		appErr := apperrors.InsertWordLibraryErr.AppendMessage(err)
		rt.log.Error(appErr)
		return appErr
	}

	return nil
}

func (rt *wordsRepository) UpdateWords(ctx context.Context, words []*models.Word) error {
	for _, word := range words {
		err := rt.UpdateWord(ctx, word)
		if err != nil {
			rt.log.Error(err)
			return err
		}
	}

	return nil
}

func (rt *wordsRepository) UpdateWord(ctx context.Context, word *models.Word) error {
	result := rt.db.Model(&models.Word{}).Where("id = ?", word.ID).
		Updates(map[string]interface{}{
			"english":         word.English,
			"russian":         word.Russian,
			"theme":           word.Theme,
			"preposition":     word.Preposition,
			"parts_of_speech": word.PartsOfSpeech,
		})
	if result.Error != nil {
		appErr := apperrors.UpdateWordErr.AppendMessage(word.English, " word ", result.Error)
		rt.log.Error(appErr)
		return appErr
	}

	if result.RowsAffected == 0 {
		appErr := &apperrors.UpdateWordRowAffectedErr
		rt.log.Info(appErr)
		return appErr
	}

	return nil
}

func (rt *wordsRepository) GetAllTopics() ([]string, error) {
	var themes []string
	err := rt.db.Table("words").Select("DISTINCT(theme)").Pluck("DISTINCT(theme)", &themes).Error
	if err != nil {
		appErr := apperrors.GetAllTopicsErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return themes, nil
}
