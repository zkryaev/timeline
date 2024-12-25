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
	LocalEnv = "LOCAL"
	// Логи выводятся на консоль
	DevEnv = "DEV"
	// Логи выводятся в файл
	ProdEnv = "PROD"
)

func New(env string) *zap.Logger {
	cfg := zap.Config{}
	encoder := zapcore.EncoderConfig{}
	filepath := envars.GetPath("APP_LOGS")
	OutputPaths := []string{"stdout"}
	if _, err := os.Stat(filepath); err == nil {
		OutputPaths = append(OutputPaths, filepath)
	} else {
		var timestamp string
		log.Println("log.txt not found")
		filename := strings.SplitAfter(filepath, "/")[len(strings.SplitAfter(filepath, "/"))-1]
		path := strings.TrimSuffix(filepath, filename)
		if env == ProdEnv {
			timestamp = time.Now().Format("15:04:05_2006-01-02_")
		}
		_, err := os.Create(path + timestamp + filename)
		if err != nil {
			log.Println("couldn't create log.txt: ", err.Error())
			os.Exit(1)
		}
		log.Println("log.txt has been created: ", path+filename)
	}
	cfg.OutputPaths = OutputPaths

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000-0700"))
	}

	switch env {
	case LocalEnv, DevEnv:
		encoder = zap.NewDevelopmentEncoderConfig()
		encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder.EncodeTime = customTimeEncoder            
		encoder.EncodeCaller = zapcore.ShortCallerEncoder

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
