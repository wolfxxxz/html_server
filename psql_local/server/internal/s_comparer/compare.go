package comparer

import (
	"net/http"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/repositories"
	"server/internal/services"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Comparer struct {
	repoUser    repositories.RepoUsers
	repoLibrary repositories.RepoLibrary
	log         *logrus.Logger
}

func NewComparer(userRepo repositories.RepoUsers, repoLibrary repositories.RepoLibrary, log *logrus.Logger) *Comparer {
	return &Comparer{repoUser: userRepo, repoLibrary: repoLibrary, log: log}
}

func (srv Comparer) CompareTestWords(r *http.Request, userID string) error {
	result := models.TestResult{}
	for i, word := range HashTableWords[userID].Words {
		answer := r.FormValue("answer" + strconv.Itoa(i))
		//srv.log.Infof("word [%v] and answer [%v]", word, answer)

		userService := services.NewUserService(srv.repoUser, srv.repoLibrary, srv.log)
		if srv.compare(word, answer) {
			//srv.log.Infof("IF COMPARE word [%v] and answer [%v]", word, answer)
			HashTableWords[userID].Words[i].Right = true
			err := userService.MoveWordToLearned(r.Context(), userID, word.ID.String())
			if err != nil {
				appErr := err.(*apperrors.AppError)
				srv.log.Error(appErr)
				return appErr
			}

			result.Right++
		} else {
			//srv.log.Infof("ELSE word [%v] and answer [%v]", word, answer)
			err := userService.AddWordToLearn(r.Context(), userID, word.ID.String())
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

func (srv Comparer) CompareLearnWords(r *http.Request, userID string) error {
	words := []*models.Word{}
	userService := services.NewUserService(srv.repoUser, srv.repoLibrary, srv.log)
	for i, word := range HashTableWordsLearn[userID].Words {
		answer := r.FormValue("answer" + strconv.Itoa(i))
		if srv.compareToLoverAndIgnoreSpace(word.English, answer) {
			err := userService.DeleteLearnFromUserById(r.Context(), userID, word.ID.String())
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
