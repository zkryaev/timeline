package logger

import (
	"fmt"
	"log"
	"os"
	"time"
	"timeline/internal/libs/envars"

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
	// получаем путь куда сохранять логи
	LogsPath := envars.GetPath("LOGS_PATH")
	OutputPaths := []string{"stdout"}
	if _, err := os.Stat(LogsPath); err == nil {
		OutputPaths = append(OutputPaths, LogsPath)
	} else {
		log.Println("Warn:", "wrong path to logs.txt")
	}
	cfg.OutputPaths = OutputPaths

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000-0700"))
	}

	switch env {
	case LocalEnv, DevEnv:
		encoder = zap.NewDevelopmentEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder.EncodeTime = customTimeEncoder            // zapcore.ISO8601TimeEncoder
		encoder.EncodeCaller = zapcore.ShortCallerEncoder // краткий вывод caller

		cfg = zap.NewDevelopmentConfig()
		cfg.DisableStacktrace = true
		cfg.EncoderConfig = encoder
	case ProdEnv:
		encoder = zap.NewProductionEncoderConfig()
		cfg = zap.NewProductionConfig()
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
