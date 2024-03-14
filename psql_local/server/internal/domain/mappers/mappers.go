package mappers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/domain/requests"
	"server/internal/domain/responses"
	"strconv"
	"unicode"

	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
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
		tempWord := &models.Word{
			ID:            libWord.ID,
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
		wordId := strconv.Itoa(word.ID)
		wordResp := &responses.WordResp{
			English:       word.English,
			Russian:       word.Russian,
			ID:            wordId,
			PartsOfSpeech: word.PartsOfSpeech,
		}

		wordsResp = append(wordsResp, wordResp)
	}

	return wordsResp
}

func MapMultipartToXLS(file *multipart.File) (*xlsx.File, error) {
	//----------------------------Map file *multipart.File-------------
	// Создаем временный файл для сохранения содержимого multipart.File
	tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		appErr := apperrors.MapMultipartToXLSErr.AppendMessage(err)
		return nil, appErr
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, *file)
	if err != nil {
		appErr := apperrors.MapMultipartToXLSErr.AppendMessage(err)
		return nil, appErr
	}

	xlFile, err := xlsx.OpenFile(tempFile.Name())
	if err != nil {
		appErr := apperrors.MapMultipartToXLSErr.AppendMessage(err)
		return nil, appErr
	}

	os.Remove(tempFile.Name())

	return xlFile, nil
}

func MapXLStoLibrary(xlFile *xlsx.File) []*models.Library {
	wordNew := []*models.Library{}
	for _, sheet := range xlFile.Sheets {
		if sheet == nil {
			break
		}

		for _, row := range sheet.Rows {
			if len(row.Cells) == 0 {
				continue
			}

			num, err := strconv.Atoi(row.Cells[0].String())
			if err != nil {
				return wordNew
			}

			word := &models.Library{
				ID:            num,
				Root:          capitalizeFirstRune(row.Cells[1].String()),
				English:       capitalizeFirstRune(row.Cells[2].String()),
				Preposition:   row.Cells[3].String(),
				Russian:       capitalizeFirstRune(row.Cells[4].String()),
				Theme:         row.Cells[5].String(),
				PartsOfSpeech: row.Cells[6].String(),
			}

			wordNew = append(wordNew, word)
		}
	}

	return wordNew
}

func MapXLStoWords(xlFile *xlsx.File) []*models.Word {
	wordNew := []*models.Word{}
	for _, sheet := range xlFile.Sheets {
		if sheet == nil {
			break
		}

		for _, row := range sheet.Rows {
			if len(row.Cells) == 0 {
				continue
			}

			num, err := strconv.Atoi(row.Cells[0].String())
			if err != nil {
				return wordNew
			}

			word := &models.Word{
				ID: num,
				//Root:          capitalizeFirstRune(row.Cells[1].String()),
				English:       capitalizeFirstRune(row.Cells[2].String()),
				Preposition:   row.Cells[3].String(),
				Russian:       capitalizeFirstRune(row.Cells[4].String()),
				Theme:         row.Cells[5].String(),
				PartsOfSpeech: row.Cells[6].String(),
			}

			wordNew = append(wordNew, word)
		}
	}

	return wordNew
}

func capitalizeFirstRune(line string) string {
	runes := []rune(line)
	for i, r := range runes {
		if i == 0 {
			runes[i] = unicode.ToUpper(r)
		}
	}

	return string(runes)
}
