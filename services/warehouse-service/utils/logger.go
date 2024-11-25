package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", "logs/warehouse-service.log"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(config.EncoderConfig)

	// Create log directory if not exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("logs/warehouse-service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), config.Level),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), config.Level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Logger = logger.Sugar()
}
