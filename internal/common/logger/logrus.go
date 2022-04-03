package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func CreateLogger(appName string) *logrus.Entry {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	return logger.WithField("app", appName)
}
