package comparer

/*
import (
	"server/internal/domain/models"
	"server/internal/repositories"
	"strings"

	"github.com/agnivade/levenshtein"
)

func (srv Comparer) compareToLoverAndIgnoreSpace(word string, answer string) bool {
	wordEnglEgnoredSpaceLoverCase := strings.ToLower(ignorSpace(word))
	answerIgnoredSpaceLoverCase := strings.ToLower(ignorSpace(answer))
	return strings.EqualFold(wordEnglEgnoredSpaceLoverCase, answerIgnoredSpaceLoverCase)
}

func (srv Comparer) compare(word *models.Word, answer string) bool {
	wordEnglEgnoredSpaceLoverCase := strings.ToLower(ignorSpace(word.English))
	answerIgnoredSpaceLoverCase := strings.ToLower(ignorSpace(answer))
	if strings.EqualFold(wordEnglEgnoredSpaceLoverCase, answerIgnoredSpaceLoverCase) {
		//srv.log.Infof("if strings.EqualFold word [%v] and answer [%v]", word, answer)
		return true
	}

	if srv.compareWithMap(word.Russian, answerIgnoredSpaceLoverCase, repositories.WordsLibraryLocalMap) {
		//srv.log.Infof("if compaRE MAP word [%v] and answer [%v]", word, answer)
		return true
	}

	return false
}

func (srv Comparer) compareStringsLevenshtein(str1, str2 string) bool {
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

func (srv Comparer) compareWithMap(russian, answerIgnoredSpaceLoverCase string, mapWords *map[string][]string) bool {
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
*/
