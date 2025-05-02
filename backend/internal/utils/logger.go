package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger создает новый экземпляр логгера
func NewLogger(level string) *logrus.Logger {
	logger := logrus.New()

	// Установка вывода в stdout
	logger.SetOutput(os.Stdout)

	// Установка формата логов
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})

	// Установка уровня логирования
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	return logger
}
