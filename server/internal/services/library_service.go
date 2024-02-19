package services

import (
	"context"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/repositories"

	"github.com/sirupsen/logrus"
)

type LibraryService struct {
	repoLibrary repositories.RepoLibrary
	log         *logrus.Logger
}

func NewLibraryService(repoLibrary repositories.RepoLibrary, log *logrus.Logger) *LibraryService {
	return &LibraryService{repoLibrary: repoLibrary, log: log}
}

func (ls *LibraryService) GetTranslationByWord(ctx context.Context, translReq *requests.TranslationRequest) ([]*models.Library, error) {
	capitalizedWord := capitalizeFirstRune(translReq.Word)
	if isCyrillic(capitalizedWord) {
		words, err := ls.repoLibrary.GetTranslationRus(capitalizedWord)
		if err != nil {
			ls.log.Error(err)
			return nil, err
		}

		if len(words) == 0 {
			words, err = ls.repoLibrary.GetTranslationRusLike(capitalizedWord)
			if err != nil {
				ls.log.Error(err)
				return nil, err
			}

		}

		return words, nil
	}

	if !isCyrillic(capitalizedWord) {
		words, err := ls.repoLibrary.GetTranslationEngl(capitalizedWord)
		if err != nil {
			ls.log.Error(err)
			return nil, err
		}

		if len(words) == 0 {
			words, err = ls.repoLibrary.GetTranslationEnglLike(capitalizedWord)
			if err != nil {
				ls.log.Error(err)
				return nil, err
			}

		}

		return words, nil
	}

	appErr := apperrors.GetTranslationByWordErr.AppendMessage("this word isn't rus or english, try to change your language")
	ls.log.Error(appErr)
	return nil, appErr
}
