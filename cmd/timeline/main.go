package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
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
		log.Fatal("No .env file found")
	}

	// Подгружаем конфиг
	cfg := config.MustLoad()

	//Инициализация логгера
	Logs := logger.New(cfg.App.Env)
	Logs.Info("Application is initializing")

	db, err := infrastructure.GetDB(os.Getenv("DB"), cfg.DB)
	if err != nil {
		Logs.Fatal("wrong db type was intered", zap.Error(err))
	}
	err = db.Open()
	if err != nil {
		Logs.Fatal(
			"failed to connect",
			zap.String("Database", os.Getenv("DB")),
			zap.Error(err),
		)
	}
	Logs.Info("Successfuly connected to", zap.String("Database", os.Getenv("DB")))
	defer db.Close()

	// Поднимаем почтовый сервис параметрами по умолчанию
	post := mail.New(cfg.Mail, Logs, 0, 0, 0)
	post.Start()
	Logs.Info("Successfuly connected to", zap.String("Mail server", os.Getenv("MAIL_HOST")))
	defer post.Shutdown()

	// Подключение к S3
	s3storage := s3.New(cfg.StorageS3)
	if err := s3storage.Connect(); err != nil {
		Logs.Fatal("failed to connect to S3", zap.Error(err))
	}
	Logs.Info("Successfully connected to S3", zap.Bool("ssl_mode", cfg.StorageS3.UseSSL))

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
	Logs.Info("Application is now running", zap.String("HTTP server", cfg.App.Host+":"+cfg.App.Port))

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
