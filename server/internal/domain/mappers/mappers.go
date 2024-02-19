package mappers

import (
	"fmt"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/domain/responses"

	"github.com/google/uuid"
)

func MapReqCreateUsToUser(userReq *requests.CreateUserRequest) *models.User {
	userID := uuid.New()
	return &models.User{
		ID:       &userID,
		Name:     userReq.Name,
		LastName: userReq.LastName,
		Email:    userReq.Email,
		Role:     userReq.Role,
	}

}

func ScanUser(u *models.User) {
	var name, password string
	fmt.Println("Your Name")
	fmt.Scan(&name)
	fmt.Println("Your Password")
	fmt.Scan(&password)
	u.Name = name
	u.Password = password
	bid := uuid.New()
	u.ID = &bid
}

func MapLibraryToWords(library []*models.Library) []*models.Word {
	words := []*models.Word{}
	for _, libWord := range library {
		byd := uuid.New()
		tempWord := &models.Word{
			ID:            &byd,
			Russian:       libWord.Russian,
			English:       libWord.English,
			Theme:         libWord.Theme,
			PartsOfSpeech: libWord.PartsOfSpeech,
		}

		words = append(words, tempWord)
	}

	return words
}

func MapLibraryToWordsGetTranslResponse(library []*models.Library) []*responses.GetTranslResponse {
	words := []*responses.GetTranslResponse{}
	for _, libWord := range library {
		tempWord := &responses.GetTranslResponse{
			Russian: libWord.Russian,
			English: libWord.English,
		}

		words = append(words, tempWord)
	}

	return words
}

func MapTokenToLoginResponse(token string, expiresAt string) *responses.LoginResponse {
	return &responses.LoginResponse{Token: token, ExpiresIn: expiresAt, TokenType: "jwt", RefreshToken: "it'll be soon"}
}

func MapWordsToWordsResp(words []*models.Word) []*responses.WordResp {
	wordsResp := []*responses.WordResp{}
	for _, word := range words {
		wordResp := &responses.WordResp{
			English:       word.English,
			Russian:       word.Russian,
			ID:            word.ID.String(),
			PartsOfSpeech: word.PartsOfSpeech,
		}

		wordsResp = append(wordsResp, wordResp)
	}

	return wordsResp
}
