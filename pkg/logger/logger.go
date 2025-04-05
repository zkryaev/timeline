package logger

import (
	"fmt"
	"log"
	"os"
	"time"
	"timeline/internal/utils/envars"
	"timeline/internal/utils/fsop"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Dev  = "DEV"
	Prod = "PROD"
)

func New(env string) *zap.Logger {
	outputPaths := []string{"stdout"}
	givenPath := envars.GetPathByEnv("APP_LOGS")
	if _, err := os.Stat(givenPath); err == nil {
		outputPaths = append(outputPaths, givenPath)
	} else {
		filepath, err := fsop.CreateDirAndFile(givenPath, true)
		if err != nil {
			log.Printf("failed to create dir/file: %s", err.Error())
		}
		log.Println("logs will be stored in: ", filepath)
		outputPaths = append(outputPaths, filepath)
	}

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000-0700"))
	}

	cfg := zap.Config{}
	encoder := zapcore.EncoderConfig{}

	switch env {
	case Dev:
		encoder = zap.NewDevelopmentEncoderConfig()
		cfg = zap.NewDevelopmentConfig()

		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder.EncodeTime = customTimeEncoder
		encoder.EncodeCaller = zapcore.ShortCallerEncoder

		cfg.OutputPaths = outputPaths
		cfg.DisableStacktrace = true
		cfg.EncoderConfig = encoder
	case Prod:
		cfg = zap.NewProductionConfig()

		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	default:
		log.Fatal("unknown env type")
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("build logger: %w", err))
	}
	return logger
}
