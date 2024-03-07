package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"server/internal/apperrors"
	"server/internal/domain/mappers"
	"server/internal/domain/models"
	"server/internal/repositories"

	"github.com/sirupsen/logrus"
)

type LibraryService struct {
	repoLibrary repositories.RepoLibrary
	repoBackup  repositories.BackUpCopyRepo
	log         *logrus.Logger
}

func NewLibraryService(repoLibrary repositories.RepoLibrary, repoBackup repositories.BackUpCopyRepo, log *logrus.Logger) *LibraryService {
	return &LibraryService{repoLibrary: repoLibrary, repoBackup: repoBackup, log: log}
}

func (ls *LibraryService) GetTranslationByWord(ctx context.Context, translReq string) ([]*models.Library, error) {
	capitalizedWord := capitalizeFirstRune(translReq)
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

func (ls *LibraryService) UpdateLibraryOldAndNewWordsByMultyFile(ctx context.Context, file *multipart.File) error {
	fileXLS, err := mappers.MapMultipartToXLS(file)
	if err != nil {
		ls.log.Error(err)
		return err
	}

	librUpdate := mappers.MapXLStoLibrary(fileXLS)

	for _, word := range librUpdate {
		err := ls.repoLibrary.UpdateWord(ctx, word)
		if err != nil {
			if err == &apperrors.UpdateWordRowAffectedErr {
				ls.log.Info(fmt.Sprintf("insert word %v, ID %v", word.English, word.ID))
				err := ls.repoLibrary.InsertWordLibrary(ctx, word)
				if err != nil {
					return err
				}

				continue
			}

			ls.log.Error(err)
			return err
		}
	}

	return nil
}

func (ls *LibraryService) DownloadXLXFromDb() (*os.File, error) {
	words, err := ls.repoLibrary.GetAllWords()
	if err != nil {
		ls.log.Error(err)
		return nil, err
	}

	ls.repoBackup.SaveWordsAsXLSX(words)
	return ls.repoBackup.OpenFile()
}

func (ls *LibraryService) GetAllTopics() ([]string, error) {
	topics, err := ls.repoLibrary.GetAllTopics()
	if err != nil {
		ls.log.Error(err)
		return nil, err
	}

	if topics == nil {
		appErr := apperrors.GetAllTopicsLibServErr.AppendMessage("topics = nil")
		ls.log.Error(appErr)
		return nil, appErr
	}

	topicsWithoutWhiteSpace := []string{}
	for _, topic := range topics {
		if topic != "" {
			topicsWithoutWhiteSpace = append(topicsWithoutWhiteSpace, topic)
		}
	}

	return topicsWithoutWhiteSpace, nil
}
