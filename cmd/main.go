package main

import (
	"os"
	"os/signal"
	"syscall"
	"timeline/internal/app"
	"timeline/internal/config"
	"timeline/internal/repository"
	"timeline/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	//Инициализация логгера
	Logr := logger.New(cfg.App.Env)

	// TODO: подключение к БД
	stubRepo := repository.StubRepo{}
	// TODO: подключение к Redis

	Application := app.New(cfg.App, Logr)
	Application.SetupControllers(cfg.Token, stubRepo)
	go Application.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	Application.Stop()
	Logr.Info(
		"Gracefully stopped",
		zap.String("signal", sig.String()),
	)

}
