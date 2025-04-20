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
// @host www.timeline.ru
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description "Bearer {token}"
// @externalDocs.description Документация
// @externalDocs.url https://github.com/zkryaev/timeline
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
	PrintConfiguration(logger, cfg)
	defer logger.Sync()
	db, err := infrastructure.GetDB(os.Getenv("DB"), &cfg.DB)
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
	defer db.Close()

	backdata := &loader.BackData{}
	if cfg.App.Settings.UseLocalBackData {
		logger.Info("Loading from local storage (DB)")
		backdata.Cities, err = db.PreLoadCities(context.Background())
		if err != nil {
			logger.Fatal("PreLoadCities", zap.Error(err))
			return
		}
	} else {
		logger.Info("Loading backdata from provided sources")
		if err := loader.LoadData(logger, db, backdata); err != nil {
			logger.Fatal("failed", zap.Error(err))
		}
	}
	logger.Info("Loading data is finished")

	var post infrastructure.Mail
	if cfg.App.Settings.EnableMail {
		// Поднимаем почтовый сервис параметрами по умолчанию
		post = mail.New(cfg.Mail, logger, 0, 0, 0)
		post.Start()
		logger.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("MAIL_HOST")))
		defer post.Shutdown()
	}

	var s3repo *s3.Minio
	if cfg.App.Settings.EnableMedia {
		// Подключение к S3
		s3repo = s3.New(cfg.S3)
		if err = s3repo.Connect(); err != nil {
			logger.Fatal(fmt.Sprintf("failed to connect to %s", os.Getenv("S3")), zap.Error(err))
		}
		logger.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("S3")))
	}
	app := app.New(cfg.App, logger)
	err = app.SetupControllers(cfg.Token, backdata, db, post, s3repo)
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
	logger.Info("Application is running")

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
	app.Shutdown(ctx)
	logger.Info("Application stopped")
}

func PrintConfiguration(logger *zap.Logger, cfg *config.Config) {
	logger.Info("Settings:")
	logger.Info("", zap.Bool("use_local_backdata", cfg.App.Settings.UseLocalBackData))

	logger.Info("Features:")
	logger.Info("", zap.Bool("enable_authorization", cfg.App.Settings.EnableAuthorization))
	logger.Info("", zap.Bool("enable_media", cfg.App.Settings.EnableMedia))
	logger.Info("", zap.Bool("enable_mail", cfg.App.Settings.EnableMail))

	logger.Info("Token's TTL:")
	logger.Info("", zap.Duration("access token", cfg.Token.AccessTTL))
	logger.Info("", zap.Duration("refresh token", cfg.Token.RefreshTTL))

	logger.Info("Server settings:")
	logger.Info("", zap.String("env-mode", cfg.App.Env))
	logger.Info("", zap.String("listening on", cfg.App.Server.Host+":"+cfg.App.Server.Port))
	logger.Info("", zap.String("request-timeout", cfg.App.Server.Timeout.String()))
	logger.Info("", zap.String("idle-timeout", cfg.App.Server.IdleTimeout.String()))

	// style formatters
	bold := "\033[1m"
	line := "\033[4m"
	reset := "\033[0m"

	logger.Info(fmt.Sprintf("Database: %s%s%s%s settings:", bold, line, cfg.DB.Protocol, reset))
	logger.Info("", zap.String("listening on", cfg.DB.Host+":"+cfg.DB.Port))
	logger.Info("", zap.String("ssl", cfg.DB.SSLmode))

	if cfg.App.Settings.EnableMail {
		logger.Info(fmt.Sprintf("Mail: %s%s%s%s settings:", bold, line, cfg.Mail.Service, reset))
		logger.Info("", zap.String("listening on", cfg.Mail.Host+":"+strconv.Itoa(cfg.Mail.Port)))
		logger.Info("", zap.String("profile", cfg.Mail.User))
	}
	if cfg.App.Settings.EnableMedia {
		logger.Info(fmt.Sprintf("S3: %s%s%s%s settings:", bold, line, cfg.S3.Name, reset))
		logger.Info("", zap.String("storage listening on", cfg.S3.Host+":"+cfg.S3.DataPort))
		logger.Info("", zap.String("console listening on", cfg.S3.Host+":"+cfg.S3.ConsolePort))
		logger.Info("", zap.Bool("ssl", cfg.S3.SSLmode))
	}
	if cfg.App.Settings.EnableMetrics {
		logger.Info(fmt.Sprintf("Metrics: %s%s%s%s settings:", bold, line, "prometheus", reset))
		logger.Info("", zap.String("listening on", cfg.Prometheus.Host+":"+cfg.Prometheus.Port))
	}
}
