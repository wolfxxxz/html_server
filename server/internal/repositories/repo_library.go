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
	InsertWordLibrary(ctx context.Context, word *models.Library) error
	UpdateWord(ctx context.Context, word *models.Library) error
	InitWordsMap() error
	GetAllTopics() ([]string, error)
}

var WordsLibraryLocalMap *map[string][]string

type repoLibrary struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewRepoLibrary(db *gorm.DB, log *logrus.Logger) RepoLibrary {
	return &repoLibrary{db: db, log: log}
}

func (rt *repoLibrary) InitWordsMap() error {
	lib, err := rt.GetAllWords()
	if err != nil {
		appErr := apperrors.InitWordsMapErr.AppendMessage(err)
		rt.log.Error(appErr)
		return appErr
	}

	wordsMap := make(map[string][]string)

	for _, word := range lib {
		wordsMap[word.Russian] = append(wordsMap[word.Russian], word.English)
	}

	WordsLibraryLocalMap = &wordsMap

	return nil
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

func (rt *repoLibrary) InsertWordLibrary(ctx context.Context, word *models.Library) error {
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

	createdLib := &models.Library{}
	if err := tx.First(createdLib, "id = ?", word.ID).Error; err != nil {
		appErr := apperrors.InsertWordLibraryErr.AppendMessage(err)
		rt.log.Error(appErr)
		return appErr
	}

	return nil
}

func (rt *repoLibrary) UpdateWord(ctx context.Context, word *models.Library) error {
	result := rt.db.Model(&models.Library{}).Where("id = ?", word.ID).
		Updates(map[string]interface{}{
			"english":         word.English,
			"russian":         word.Russian,
			"theme":           word.Theme,
			"preposition":     word.Preposition,
			"parts_of_speech": word.PartsOfSpeech,
			"root":            word.Root,
		})
	if result.Error != nil {
		appErr := apperrors.UpdateWordErr.AppendMessage(result.Error)
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

func (rt *repoLibrary) GetAllTopics() ([]string, error) {
	var themes []string
	err := rt.db.Table("libraries").Select("DISTINCT(theme)").Pluck("DISTINCT(theme)", &themes).Error
	if err != nil {
		appErr := apperrors.GetAllTopicsErr.AppendMessage(err)
		rt.log.Error(appErr)
		return nil, appErr
	}

	return themes, nil
}
