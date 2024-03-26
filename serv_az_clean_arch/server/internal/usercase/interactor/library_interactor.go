package interactor

import (
	"context"
	"mime/multipart"
	"os"
	"server/internal/apperrors"
	"server/internal/domain/mappers"
	"server/internal/domain/models"
	"server/internal/usercase/repository"
)

type libraryInteractor struct {
	LibraryRepository repository.LibraryRepository
	WordsRepository   repository.WordsRepository
	BackupRepository  repository.BackUpCopyRepo
}

type LibraryInteractor interface {
	GetTranslationByWord(ctx context.Context, translReq string) ([]*models.Library, error)
	GetTranslationByPieceOfWord(ctx context.Context, translReq string) (string, error)
	UpdateLibraryOldAndNewWordsByMultyFile(ctx context.Context, file *multipart.File) error
	DownloadXLXFromDb() (*os.File, error)
	GetAllTopics() ([]string, error)
}

func NewLibraryInteractor(u repository.LibraryRepository, w repository.WordsRepository, b repository.BackUpCopyRepo) LibraryInteractor {
	return &libraryInteractor{LibraryRepository: u, WordsRepository: w, BackupRepository: b}
}

func (ls *libraryInteractor) GetTranslationByWord(ctx context.Context, translReq string) ([]*models.Library, error) {
	capitalizedWord := capitalizeFirstRune(translReq)
	if isCyrillic(capitalizedWord) {
		words, err := ls.LibraryRepository.GetTranslationRus(capitalizedWord)
		if err != nil {
			return nil, err
		}

		if len(words) == 0 {
			words, err = ls.LibraryRepository.GetTranslationRusLike(capitalizedWord)
			if err != nil {
				return nil, err
			}

		}

		return words, nil
	}

	if !isCyrillic(capitalizedWord) {
		words, err := ls.LibraryRepository.GetTranslationEngl(capitalizedWord)
		if err != nil {
			return nil, err
		}

		if len(words) == 0 {
			words, err = ls.LibraryRepository.GetTranslationEnglLike(capitalizedWord)
			if err != nil {
				return nil, err
			}

		}

		return words, nil
	}

	appErr := apperrors.GetTranslationByWordErr.AppendMessage("this word isn't rus or english, try to change your language")

	return nil, appErr
}

func (ls *libraryInteractor) GetTranslationByPieceOfWord(ctx context.Context, translReq string) (string, error) {
	capitalizedWord := capitalizeFirstRune(translReq)
	if isCyrillic(capitalizedWord) {
		words, err := ls.LibraryRepository.GetTranslationRusLikeWord(capitalizedWord)
		if err != nil {
			//ls.log.Error(err)
			return "", err
		}

		return words.Russian, nil
	}

	if !isCyrillic(capitalizedWord) {
		words, err := ls.LibraryRepository.GetTranslationEnglLikeWord(capitalizedWord)
		if err != nil {
			//ls.log.Error(err)
			return "", err
		}

		return words.English, nil
	}

	appErr := apperrors.GetTranslationByWordErr.AppendMessage("this word isn't rus or english, try to change your language")
	//ls.log.Error(appErr)
	return "", appErr
}

func (ls *libraryInteractor) UpdateLibraryOldAndNewWordsByMultyFile(ctx context.Context, file *multipart.File) error {
	fileXLS, err := mappers.MapMultipartToXLS(file)
	if err != nil {
		return err
	}

	librUpdate := mappers.MapXLStoLibrary(fileXLS)

	for _, word := range librUpdate {
		err := ls.LibraryRepository.UpdateWord(ctx, word)
		if err != nil {
			if err == &apperrors.UpdateWordRowAffectedErr {
				err := ls.LibraryRepository.InsertWordLibrary(ctx, word)
				if err != nil {
					return err
				}

				continue
			}

			return err
		}
	}
	//
	err = ls.LibraryRepository.UpdateWordsMap()
	if err != nil {
		return err
	}

	wordsUpdate := mappers.MapXLStoWords(fileXLS)
	for _, word := range wordsUpdate {
		err := ls.WordsRepository.UpdateWord(ctx, word)
		if err != nil {
			if err == &apperrors.UpdateWordRowAffectedErr {
				err := ls.WordsRepository.InsertWord(ctx, word)
				if err != nil {
					return err
				}

				continue
			}

			return err
		}
	}

	return nil
}

func (ls *libraryInteractor) DownloadXLXFromDb() (*os.File, error) {
	words, err := ls.LibraryRepository.GetAllWords()
	if err != nil {
		return nil, err
	}

	ls.BackupRepository.SaveWordsAsXLSX(words)
	return ls.BackupRepository.OpenFile()
}

func (ls *libraryInteractor) GetAllTopics() ([]string, error) {
	topics, err := ls.LibraryRepository.GetAllTopics()
	if err != nil {
		return nil, err
	}

	if topics == nil {
		appErr := apperrors.GetAllTopicsLibServErr.AppendMessage("topics = nil")
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
