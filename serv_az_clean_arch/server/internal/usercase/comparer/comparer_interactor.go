package comparer

import (
	"net/http"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/interface/repository"
	"server/internal/usercase/interactor"
	"strconv"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/sirupsen/logrus"
)

// i think this is a bad way to save some data, it's better to use reddis or some hashDB
var HashTableWords = make(map[string]*models.TestPageData)
var HashTableWordsLearn = make(map[string]*models.TestPageData)

type Comparer interface {
	CompareTestWords(r *http.Request, userID string) error
	CompareLearnWords(r *http.Request, userID string) error
}

type comparer struct {
	LibraryInteractor interactor.LibraryInteractor
	UserInteractor    interactor.UserInteractor
	log               *logrus.Logger
}

func NewComparer(LibraryInteractor interactor.LibraryInteractor,
	UserInteractor interactor.UserInteractor, log *logrus.Logger) Comparer {
	return &comparer{
		LibraryInteractor: LibraryInteractor,
		UserInteractor:    UserInteractor,
		log:               log,
	}
}

func (srv comparer) CompareTestWords(r *http.Request, userID string) error {
	result := models.TestResult{}
	for i, word := range HashTableWords[userID].Words {
		answer := r.FormValue("answer" + strconv.Itoa(i))
		//srv.log.Infof("word [%v] and answer [%v]", word, answer)

		wordId := strconv.Itoa(word.ID)
		if srv.compare(word, answer) {
			//srv.log.Infof("IF COMPARE word [%v] and answer [%v]", word, answer)
			HashTableWords[userID].Words[i].Right = true

			err := srv.UserInteractor.MoveWordToLearned(r.Context(), userID, wordId)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.log.Error(appErr)
				return appErr
			}

			result.Right++
		} else {
			//srv.log.Infof("ELSE word [%v] and answer [%v]", word, answer)
			err := srv.UserInteractor.AddWordToLearn(r.Context(), userID, wordId)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.log.Error(appErr)
				return appErr
			}

			result.Wrong++
		}

	}

	HashTableWords[userID].Result = &result
	HashTableWords[userID].TestPassed = true

	return nil
}

func (srv comparer) CompareLearnWords(r *http.Request, userID string) error {
	words := []*models.Word{}
	for i, word := range HashTableWordsLearn[userID].Words {
		answer := r.FormValue("answer" + strconv.Itoa(i))
		if srv.compareToLoverAndIgnoreSpace(word.English, answer) {
			wordId := strconv.Itoa(word.ID)
			err := srv.UserInteractor.DeleteLearnFromUserById(r.Context(), userID, wordId)
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.log.Error(appErr)
				return appErr
			}
		} else {
			words = append(words, word)
		}
	}

	HashTableWordsLearn[userID].Words = words

	if len(words) == 0 {
		HashTableWordsLearn[userID].LearnPassed = true
	}

	return nil
}

func (srv comparer) compareToLoverAndIgnoreSpace(word string, answer string) bool {
	wordEnglEgnoredSpaceLoverCase := strings.ToLower(ignorSpace(word))
	answerIgnoredSpaceLoverCase := strings.ToLower(ignorSpace(answer))
	return strings.EqualFold(wordEnglEgnoredSpaceLoverCase, answerIgnoredSpaceLoverCase)
}

func (srv comparer) compare(word *models.Word, answer string) bool {
	wordEnglEgnoredSpaceLoverCase := strings.ToLower(ignorSpace(word.English))
	answerIgnoredSpaceLoverCase := strings.ToLower(ignorSpace(answer))
	if strings.EqualFold(wordEnglEgnoredSpaceLoverCase, answerIgnoredSpaceLoverCase) {
		//srv.log.Infof("if strings.EqualFold word [%v] and answer [%v]", word, answer)
		return true
	}

	if srv.compareWithMap(word.Russian, answerIgnoredSpaceLoverCase, repository.WordsLibraryLocalMap) {
		//srv.log.Infof("if compaRE MAP word [%v] and answer [%v]", word, answer)
		return true
	}

	return false
}

func (srv comparer) compareStringsLevenshtein(str1, str2 string) bool {
	str1 = ignorSpace(strings.ToLower(str1))
	str2 = ignorSpace(strings.ToLower(str2))
	mistakes := 1
	if distance := levenshtein.ComputeDistance(str1, str2); distance <= mistakes {
		return true
	}

	return false
}

func ignorSpace(s string) (c string) {
	for _, v := range s {
		if v != ' ' {
			c = c + string(v)
		}
	}

	return
}

func (srv comparer) compareWithMap(russian, answerIgnoredSpaceLoverCase string, mapWords *map[string][]string) bool {
	englishWords, ok := (*mapWords)[russian]
	if !ok {
		return false
	}

	for _, word := range englishWords {
		if srv.compareStringsLevenshtein(answerIgnoredSpaceLoverCase, word) {
			return true
		}
	}

	return false
}
