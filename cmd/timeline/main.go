package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"timeline/internal/app"
	"timeline/internal/config"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/mail"
	"timeline/internal/infrastructure/s3"
	"timeline/internal/sugar/cronjob"
	"timeline/internal/utils/loader"
	"timeline/pkg/logger"

	"github.com/joho/godotenv"
	_ "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Timeline API
// @version 1.0
// @BasePath /v1
// @schemes http
func main() {
	// Подгружаем все переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not found!")
	}

	// Подгружаем конфиг
	cfg := config.MustLoad()
	successConnection := "Successfuly connected to"
	// Инициализация логгера
	logger := logger.New(cfg.App.Env)
	logger.Info("Application started")
	defer logger.Sync()
	db, err := infrastructure.GetDB(os.Getenv("DB"), cfg.DB)
	if err != nil {
		logger.Fatal("incorrect db type", zap.Error(err))
	}
	err = db.Open()
	if err != nil {
		logger.Fatal(
			fmt.Sprintf("failed to connect %s", os.Getenv("DB")),
			zap.Error(err),
		)
	}
	logger.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("DB")))
	logger.Info(
		"DB",
		zap.String("database server", cfg.DB.Host+":"+cfg.DB.Port),
		zap.String("ssl", cfg.DB.SSLmode),
	)
	defer db.Close()

	if !cfg.App.IsBackDataLoaded {
		logger.Info("Loading data from provided sources...")
		if err := loader.LoadData(logger, db); err != nil {
			logger.Fatal("failed", zap.Error(err))
		}
		logger.Info("Loading data is finished")
	} else {
		logger.Info("Skipped loading background data from sources")
	}

	// Поднимаем почтовый сервис параметрами по умолчанию
	post := mail.New(cfg.Mail, logger, 0, 0, 0)
	post.Start()
	logger.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("MAIL_HOST")))
	logger.Info(
		"MAIL",
		zap.String("mail server", cfg.Mail.Host+":"+strconv.Itoa(cfg.Mail.Port)),
	)
	defer post.Shutdown()

	// Подключение к S3
	s3storage := s3.New(cfg.S3)
	if err = s3storage.Connect(); err != nil {
		logger.Fatal(fmt.Sprintf("failed to connect to %s", os.Getenv("S3")), zap.Error(err))
	}
	logger.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("S3")))
	logger.Info(
		"S3",
		zap.String("storage", cfg.S3.Host+":"+cfg.S3.DataPort),
		zap.String("console", cfg.S3.Host+":"+cfg.S3.ConsolePort),
		zap.Bool("ssl", cfg.S3.SSLmode),
	)

	app := app.New(cfg.App, logger)
	err = app.SetupControllers(cfg.Token, db, post, s3storage)
	if err != nil {
		logger.Fatal(
			"failed setup controllers",
			zap.Error(err),
		)
	}

	s := cronjob.InitCronScheduler(db)
	defer s.Shutdown()
	s.Start()

	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	errorChan := make(chan error, 1)
	go func() {
		err = app.Run()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Error("failed to run server", zap.Error(err))
				errorChan <- err
			}
		}
	}()
	logger.Info("application is running")
	logger.Info("", zap.String("app server", cfg.App.Host+":"+cfg.App.Port))

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		cancel()
		logger.Info("Received signal",
			zap.String("signal", sig.String()),
		)
	case err = <-errorChan:
		cancel()
		logger.Error("error occurred",
			zap.Error(err),
		)
	}
	app.Stop(ctx)
	logger.Info("Application stopped")
}
