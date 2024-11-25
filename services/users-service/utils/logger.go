package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	if Logger != nil {
		return // Skip if logger already initialized (for tests)
	}

	// Buat direktori logs jika belum ada
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	logFile := filepath.Join(logDir, "app.log")

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logFile}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func LogError(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, fields...)
}

func LogInfo(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, fields...)
}

func LogDebug(msg string, fields ...zapcore.Field) {
	Logger.Debug(msg, fields...)
}
