package log

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Level string
}

func InitLogrus(cfg Config) {
	logLevel, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
}
