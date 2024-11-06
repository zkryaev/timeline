package logger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LocalEnv = "LOCAL"
	DevEnv   = "DEV"
	ProdEnv  = "PROD"
)

func New(env string) *zap.Logger {
	cfg := zap.Config{}
	encoder := zapcore.EncoderConfig{}
	switch env {
	case LocalEnv, DevEnv:
		encoder = zap.NewDevelopmentEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder.EncodeTime = zapcore.ISO8601TimeEncoder

		cfg = zap.NewDevelopmentConfig()
		cfg.OutputPaths = []string{"stdout", "logs/logs.txt"}
		cfg.DisableStacktrace = true
		cfg.EncoderConfig = encoder
	case ProdEnv:
		encoder = zap.NewProductionEncoderConfig()
		cfg = zap.NewProductionConfig()
		// TODO: сделать чтобы для прода, путь задавался через переменую окружения, либо сделать через передачу аргумента из конфига приложения
		cfg.EncoderConfig = encoder
	case "":
		log.Fatal("logger did't setup: ENV is empty")
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("build logger: %w", err))
	}
	return logger
}
