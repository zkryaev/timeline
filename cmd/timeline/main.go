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
	"timeline/internal/libs/cronjob"
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
	logs := logger.New(cfg.App.Env)
	logs.Info("application initializing...")
	defer logs.Sync()
	db, err := infrastructure.GetDB(os.Getenv("DB"), cfg.DB)
	if err != nil {
		logs.Fatal("incorrect db type", zap.Error(err))
	}
	err = db.Open()
	if err != nil {
		logs.Fatal(
			fmt.Sprintf("failed to connect %s", os.Getenv("DB")),
			zap.Error(err),
		)
	}
	logs.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("DB")))
	logs.Info(
		"",
		zap.String("database server", cfg.DB.Host+":"+cfg.DB.Port),
		zap.String("ssl", cfg.DB.SSLmode),
	)
	defer db.Close()

	// Поднимаем почтовый сервис параметрами по умолчанию
	post := mail.New(cfg.Mail, logs, 0, 0, 0)
	post.Start()
	logs.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("MAIL_HOST")))
	logs.Info(
		"",
		zap.String("mail server", cfg.Mail.Host+":"+strconv.Itoa(cfg.Mail.Port)),
	)
	defer post.Shutdown()

	// Подключение к S3
	s3storage := s3.New(cfg.S3)
	if err = s3storage.Connect(); err != nil {
		logs.Fatal(fmt.Sprintf("failed to connect to %s", os.Getenv("S3")), zap.Error(err))
	}
	logs.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("S3")))
	logs.Info(
		"",
		zap.String("storage", cfg.S3.Host+":"+cfg.S3.DataPort),
		zap.String("console", cfg.S3.Host+":"+cfg.S3.ConsolePort),
		zap.Bool("ssl", cfg.S3.SSLmode),
	)

	app := app.New(cfg.App, logs)
	err = app.SetupControllers(cfg.Token, db, post, s3storage)
	if err != nil {
		logs.Fatal(
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
				logs.Error("failed to run server", zap.Error(err))
				errorChan <- err
			}
		}
	}()
	logs.Info("application is running")
	logs.Info("", zap.String("app server", cfg.App.Host+":"+cfg.App.Port))

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		cancel()
		logs.Info("Received signal",
			zap.String("signal", sig.String()),
		)
	case err = <-errorChan:
		cancel()
		logs.Error("error occurred",
			zap.Error(err),
		)
	}
	app.Stop(ctx)
	logs.Info("application stopped")
}
