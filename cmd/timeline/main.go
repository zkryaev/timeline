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

	_ "github.com/jackc/pgx/v5/stdlib"
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
	//Инициализация логгера
	Logs := logger.New(cfg.App.Env)
	Logs.Info("Application initializing...")
	defer Logs.Sync()
	db, err := infrastructure.GetDB(os.Getenv("DB"), cfg.DB)
	if err != nil {
		Logs.Fatal("incorrect db type", zap.Error(err))
	}
	err = db.Open()
	if err != nil {
		Logs.Fatal(
			fmt.Sprintf("failed to connect %s", os.Getenv("DB")),
			zap.Error(err),
		)
	}
	Logs.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("DB")))
	Logs.Info(
		"",
		zap.String("database server", cfg.DB.Host+":"+cfg.DB.Port),
		zap.String("ssl", cfg.DB.SSLmode),
	)
	defer db.Close()

	// Поднимаем почтовый сервис параметрами по умолчанию
	post := mail.New(cfg.Mail, Logs, 0, 0, 0)
	post.Start()
	Logs.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("MAIL_HOST")))
	Logs.Info(
		"",
		zap.String("mail server", cfg.Mail.Host+":"+strconv.Itoa(cfg.Mail.Port)),
	)
	defer post.Shutdown()

	// Подключение к S3
	s3storage := s3.New(cfg.S3)
	if err := s3storage.Connect(); err != nil {
		Logs.Fatal(fmt.Sprintf("failed to connect to %s", os.Getenv("S3")), zap.Error(err))
	}
	Logs.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("S3")))
	Logs.Info(
		"",
		zap.String("storage", cfg.S3.Host+":"+cfg.S3.DataPort),
		zap.String("console", cfg.S3.Host+":"+cfg.S3.ConsolePort),
		zap.Bool("ssl", cfg.S3.SSLmode),
	)

	App := app.New(cfg.App, Logs)
	err = App.SetupControllers(cfg.Token, db, post, s3storage)
	if err != nil {
		Logs.Fatal(
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
		err := App.Run()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				Logs.Error("failed to run server", zap.Error(err))
				errorChan <- err
			}
		}
	}()
	Logs.Info("Application is running")
	Logs.Info("", zap.String("app server", cfg.App.Host+":"+cfg.App.Port))

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		cancel()
		Logs.Info("Received signal",
			zap.String("signal", sig.String()),
		)
	case err := <-errorChan:
		cancel()
		Logs.Error("error occurred",
			zap.Error(err),
		)
	}
	App.Stop(ctx)
	Logs.Info("Application stopped")
}
