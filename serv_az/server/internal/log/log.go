package log

import (
	"os"
	"server/internal/apperrors"

	"github.com/sirupsen/logrus"
)

func NewLogAndSetLevel(logLevel string) (*logrus.Logger, error) {
	log := logrus.New()
	loggerLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		appErr := apperrors.NewLoggerErr.AppendMessage(err)
		return nil, appErr
	}

	log.SetLevel(loggerLevel)
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.Info("Logger has been configurated")
	return log, nil
}

func SetLevel(log *logrus.Logger, logLevel string) error {
	loggerLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		appErr := apperrors.SetLevelErr.AppendMessage(err)
		return appErr
	}

	log.SetLevel(loggerLevel)
	log.Info("logger level has been configurated")
	return nil
}
