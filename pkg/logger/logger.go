package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"timeline/internal/libs/envars"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// Логи выводятся на консоль
	Local = "LOCAL"
	// Логи выводятся на консоль
	Dev = "DEV"
	// Логи выводятся в файл
	Prod = "PROD"
)

func New(env string) *zap.Logger {
	filepath := envars.GetPathByEnv("APP_LOGS")
	outputPaths := []string{"stdout"}
	if _, err := os.Stat(filepath); err == nil {
		outputPaths = append(outputPaths, filepath)
	} else {
		var timestamp string
		log.Println("log.txt not found")
		filename := strings.SplitAfter(filepath, "/")[len(strings.SplitAfter(filepath, "/"))-1]
		path := strings.TrimSuffix(filepath, filename)
		if env == Prod {
			timestamp = time.Now().Format("15:04:05_2006-01-02_")
		}
		_, err = os.Create(path + timestamp + filename)
		if err != nil {
			log.Println("couldn't create log.txt: ", err.Error())
			os.Exit(1)
		}
		log.Println("log.txt has been created: ", path+filename)
	}

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000-0700"))
	}

	cfg := zap.Config{}
	encoder := zapcore.EncoderConfig{}

	switch env {
	case Local, Dev:
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
	case "":
		log.Fatal("logger did't setup: ENV is empty")
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("build logger: %w", err))
	}
	return logger
}
