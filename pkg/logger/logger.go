package logger

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
	"timeline/internal/libs/envars"

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
		pathparts := strings.SplitAfter(givenPath, "/")
		filename := pathparts[len(strings.SplitAfter(givenPath, "/"))-1]
		timestamp := time.Now().Format("15:04:05_2006-01-02_")
		pathDir := strings.TrimSuffix(givenPath, filename)
		if _, err := os.Stat(pathDir); errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(pathDir, os.ModePerm); err != nil {
				log.Fatalln("couldn't create logs dir: ", err.Error())
			}
		}
		filepath := pathDir + timestamp + filename
		if _, err = os.Create(filepath); err != nil && errors.Is(err, fs.ErrExist) {
			log.Fatalln("couldn't create log.txt: ", err.Error())
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
