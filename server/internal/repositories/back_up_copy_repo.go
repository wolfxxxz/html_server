package repositories

import (
	"encoding/json"
	"io"
	"os"
	"server/internal/apperrors"
	"server/internal/domain/models"

	"github.com/sirupsen/logrus"
)

type BackUpCopyRepo interface {
	GetAllFromBackUp() ([]*models.Library, error)
	SaveAllAsJson(s []*models.Library) error
}

type backUpCopyRepo struct {
	reserveCopyPath    string
	reserveCopyPathTXT string
	log                *logrus.Logger
}

func NewBackUpCopyRepo(path string, pathTXT string, log *logrus.Logger) BackUpCopyRepo {
	return &backUpCopyRepo{reserveCopyPath: path, reserveCopyPathTXT: pathTXT, log: log}
}

func (tr *backUpCopyRepo) GetAllFromBackUp() ([]*models.Library, error) {
	filejson, err := os.Open(tr.reserveCopyPath)
	if err != nil {
		appErr := apperrors.GetAllFromBackUpErr.AppendMessage(err)
		tr.log.Error(appErr)
		return nil, appErr
	}

	defer filejson.Close()
	data, err := io.ReadAll(filejson)
	if err != nil {
		appErr := apperrors.GetAllFromBackUpErr.AppendMessage(err)
		tr.log.Error(appErr)
		return nil, appErr
	}

	var words []*models.Library
	err = json.Unmarshal(data, &words)
	if err != nil {
		appErr := apperrors.GetAllFromBackUpErr.AppendMessage(err)
		tr.log.Error(appErr)
		return nil, appErr
	}

	return words, nil
}

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
