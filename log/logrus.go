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

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		TimestampFormat:        time.RFC3339,
		FullTimestamp:          true,
	})
}
