package repository

import (
	"os"
	"server/internal/apperrors"
	"server/internal/domain/models"
	"server/internal/usercase/repository"

	"strconv"
	"unicode"

	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

type backUpCopyRepo struct {
	copyPathXLSX string
	log          *logrus.Logger
}

func NewBackUpCopyRepo(copyPathXLSX string, log *logrus.Logger) repository.BackUpCopyRepo {
	return &backUpCopyRepo{copyPathXLSX: copyPathXLSX, log: log}
}

func (tr *backUpCopyRepo) GetAllWordsFromBackUpXlsx() ([]*models.Library, error) {
	xlFile, err := xlsx.OpenFile(tr.copyPathXLSX)
	if err != nil {
		appErr := apperrors.GetAllWordsFromBackUpXlsxErr.AppendMessage(err)
		tr.log.Error(appErr)
		return nil, appErr
	}

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
				//appErr := apperrors.GetAllWordsFromBackUpXlsxErr.AppendMessage(err)
				tr.log.Info("wrong created file.xlsx")
				tr.log.Infof("[%v] words has been added right now", len(wordNew))
				return wordNew, nil
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

			tr.log.Info(word.English)
		}
	}

	tr.log.Infof("[%v] words has been added", len(wordNew))
	return wordNew, nil
}

func (tr *backUpCopyRepo) SaveWordsAsXLSX(words []*models.Library) error {
	file := xlsx.NewFile()

	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		appErr := apperrors.SaveWordsAsXLSXErr.AppendMessage(err)
		tr.log.Error(err)
		return appErr
	}

	for _, word := range words {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.SetInt(word.ID)
		cell = row.AddCell()
		cell.Value = word.Root
		cell = row.AddCell()
		cell.Value = word.English
		cell = row.AddCell()
		cell.Value = word.Preposition
		cell = row.AddCell()
		cell.Value = word.Russian
		cell = row.AddCell()
		cell.Value = word.Theme
		cell = row.AddCell()
		cell.Value = word.PartsOfSpeech
		//cell = row.AddCell()
		//cell.SetInt(word.RightAnswer)
	}

	err = file.Save(tr.copyPathXLSX)
	if err != nil {
		appErr := apperrors.SaveWordsAsXLSXErr.AppendMessage(err)
		tr.log.Error(err)
		return appErr
	}

	return nil
}

func (tr *backUpCopyRepo) OpenFile() (*os.File, error) {
	file, err := os.Open(tr.copyPathXLSX)
	if err != nil {
		tr.log.Error(err)
		return nil, err
	}

	return file, nil
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

/*

func (tr *backUpCopyRepo) SaveAllAsJson(s []*models.Library) error {
	byteArr, err := json.MarshalIndent(s, "", "   ")
	if err != nil {
		appErr := apperrors.SaveAllAsJsonErr.AppendMessage(err)
		tr.log.Error(appErr)
		return err
	}

	err = os.WriteFile(tr.reserveCopyPath, byteArr, 0666) //-rw-rw-rw- 0664
	if err != nil {
		appErr := apperrors.SaveAllAsJsonErr.AppendMessage(err)
		tr.log.Error(appErr)
		return err
	}

	return nil
}
*/
