package log

import (
	"io"
	"os"
	"server/internal/apperrors"

	"github.com/sirupsen/logrus"
)

const loggerFile = "log/server.log"

func NewLogAndSetLevel(logLevel string) (*logrus.Logger, error) {
	log := logrus.New()
	loggerLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		appErr := apperrors.NewLoggerErr.AppendMessage(err)
		return nil, appErr
	}

	log.SetLevel(loggerLevel)
	log.SetReportCaller(true)
	//set output
	file, err := os.OpenFile(loggerFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
	}

	log.Infof("file: [%v] opened success\n", loggerFile)
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
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
