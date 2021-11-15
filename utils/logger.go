package utils

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"time"
)

func InitLogger() {
	logrus.Infoln("init logger")
	var logLevel = logrus.InfoLevel
	if Cfg.Debug {
		logLevel = logrus.DebugLevel
	}
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
	logrus.AddHook(rotateFileHook)
}
