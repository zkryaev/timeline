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
	_ "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Timeline API
// @version 1.0
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	cfg := config.MustLoad()

	//Инициализация логгера
	Logr := logger.New(cfg.App.Env)
	// TODO: Надо бы Авто миграцию
	db := postgres.New(cfg.DB)
	err := db.Open()
	if err != nil {
		Logr.Fatal(
			"failed connection to Database",
			zap.Error(err),
		)
	}
	defer db.Close()
	// TODO: Redis

	// Поднимаем почтовый сервис
	mail := notify.New(cfg.Mail)

	Application := app.New(cfg.App, Logr)
	err = Application.SetupControllers(cfg.Token, db, mail)
	if err != nil {
		Logr.Fatal(
			"failed setup controllers",
			zap.Error(err),
		)
	}

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
