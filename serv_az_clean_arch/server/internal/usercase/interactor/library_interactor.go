package interactor

import (
	"context"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/usercase/repository"
)

type libraryInteractor struct {
	LibraryRepository repository.LibraryRepository
}

type LibraryInteractor interface {
	GetTranslationByWord(ctx context.Context, translReq string) ([]*models.Library, error)
	GetTranslationByPieceOfWord(ctx context.Context, translReq string) (string, error)
}

func NewLibraryInteractor(u repository.LibraryRepository) LibraryInteractor {
	return &libraryInteractor{LibraryRepository: u}
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

/*
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
	//
	err = ls.repoLibrary.UpdateWordsMap()
	if err != nil {
		ls.log.Error(err)
		return err
	}

	wordsUpdate := mappers.MapXLStoWords(fileXLS)
	for _, word := range wordsUpdate {
		err := ls.repoWords.UpdateWord(ctx, word)
		if err != nil {
			if err == &apperrors.UpdateWordRowAffectedErr {
				ls.log.Info(fmt.Sprintf("insert word %v, ID %v", word.English, word.ID))
				err := ls.repoWords.InsertWord(ctx, word)
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
*/
