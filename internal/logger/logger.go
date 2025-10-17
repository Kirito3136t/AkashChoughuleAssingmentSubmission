package logger

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			Log.Errorf("Failed to create log directory: %v", err)
		}
	}

	logFileName := logDir + "/logs_" + time.Now().Format("20060102_150405") + ".log"
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		Log.Errorf("Failed to open log file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	Log.SetOutput(mw)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	Log.SetLevel(logrus.InfoLevel)

	Log.Info("Logger initialized")
}
