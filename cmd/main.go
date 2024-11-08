package main

import (
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
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	cfg := config.MustLoad()

	//Инициализация логгера
	Logr := logger.New(cfg.App.Env)

	// TODO: подключение к БД
	repo := &postgres.PostgresRepo{}
	// TODO: подключение к Redis

	// Поднимаем почтовый сервис
	mail := notify.New(cfg.Mail)

	Application := app.New(cfg.App, Logr)
	err := Application.SetupControllers(cfg.Token, repo, mail)
	if err != nil {
		Logr.Fatal(
			"failed setup controllers",
			zap.Error(err),
		)
	}

	// TODO: обсудить с Захаром легально так делать
	quit := make(chan os.Signal, 1)
	go func() {
		err = Application.Run()
		if err != nil {
			quit <- os.Interrupt
		}
	}()

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	if err != nil {
		Logr.Fatal(
			"failed to run server",
			zap.Error(err),
		)
	}
	Application.Stop()
	Logr.Info(
		"Gracefully stopped",
		zap.String("signal", sig.String()),
	)

}
