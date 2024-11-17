package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"timeline/internal/app"
	"timeline/internal/config"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mail/notify"
	"timeline/pkg/logger"

	"github.com/joho/godotenv"
	_ "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Timeline API
// @version 1.0
func main() {
	// Подгружаем все переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// Подгружаем конфиг
	cfg := config.MustLoad()

	//Инициализация логгера
	Logs := logger.New(cfg.App.Env)
	Logs.Info("Application is launched")

	db := postgres.New(cfg.DB)
	err := db.Open()
	if err != nil {
		Logs.Fatal(
			"failed connection to Database",
			zap.Error(err),
		)
	}
	Logs.Info("Connected to Database successfuly")
	defer db.Close()
	// TODO: Redis

	// Поднимаем почтовый сервис
	mail := notify.New(cfg.Mail)
	Logs.Info("Connected to Mail server successfuly")

	App := app.New(cfg.App, Logs)
	err = App.SetupControllers(cfg.Token, db, mail)
	if err != nil {
		Logs.Fatal(
			"failed setup controllers",
			zap.Error(err),
		)
	}

	quit := make(chan os.Signal, 1)
	errorChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := App.Run()
		if err != nil {
			Logs.Error("failed to run server", zap.Error(err))
		}
		errorChan <- err
	}()
	Logs.Info("Application is started")

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
