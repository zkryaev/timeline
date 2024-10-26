package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"timeline/internal/config"
	"timeline/internal/controller"
	"timeline/internal/custom"

	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	// Инициализация логгера
	custom.SetupLogger(cfg.App.Env)
	custom.Logger.Info("Configuration loaded successfully")

	// подключение к БД

	// Инициализация запуск сервера
	srv := &http.Server{
		Addr:         cfg.App.Host + cfg.App.Port,
		Handler:      controller.InitRouter(), //,
		ReadTimeout:  cfg.App.Timeout,
		WriteTimeout: cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil {
			custom.Logger.Error(
				"failed to start server",
				zap.Error(err),
			)
		}
	}()
	wg.Wait()

	// gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	custom.Logger.Info(
		"Received shutdown signal",
		zap.String("signal", sig.String()),
	)

}
