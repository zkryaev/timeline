package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"timeline/internal/app"
	"timeline/internal/config"
	"timeline/internal/controller/monitoring/metrics"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not found!")
	}

	// Подгружаем конфиг
	cfg := config.MustLoad()
	successConnection := "Successfuly connected to"
	// Инициализация логгера
	logger := logger.New(cfg.App.Stage)
	defer logger.Sync()

	logger.Info("Application started")
	PrintConfiguration(logger, cfg)

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

	s := cronjob.InitCronScheduler(db)
	defer s.Shutdown()
	s.Start()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer logger.Info("All launched services is closed")
	timeout := 1 * time.Minute
	errch := make(chan error, 2)

	app := app.New(cfg.App, logger)
	err = app.SetupControllers(cfg, backdata, db, post, s3repo)
	if err != nil {
		logger.Fatal(
			"failed setup controllers",
			zap.Error(err),
		)
	}
	app.Run(errch)
	defer app.Shutdown(ctx, timeout)
	logger.Info("Application server is listening")

	var promHandler *metrics.Prometheus
	if cfg.App.Settings.EnableMetrics {
		promHandler = metrics.NewPrometheusHandler(cfg.Prometheus, logger)
		promHandler.Launch(errch)
		defer promHandler.Shutdown(ctx, timeout)
		logger.Info("Prometheus handler is listening")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		logger.Info("Received signal", zap.String("signal", sig.String()))
	case err = <-errch:
		logger.Error("error occurred", zap.Error(err))
		cancel()
	}
}

func PrintConfiguration(logger *zap.Logger, cfg *config.Config) {
	logger.Info("Settings:")
	logger.Info("", zap.Bool("use_local_backdata", cfg.App.Settings.UseLocalBackData))

	logger.Info("Features:")
	logger.Info("", zap.Bool("enable_authorization", cfg.App.Settings.EnableAuthorization))
	logger.Info("", zap.Bool("enable_media", cfg.App.Settings.EnableMedia))
	logger.Info("", zap.Bool("enable_mail", cfg.App.Settings.EnableMail))
	logger.Info("", zap.Bool("enable_metrics", cfg.App.Settings.EnableMetrics))

	logger.Info("Token's TTL:")
	logger.Info("", zap.Duration("access token", cfg.Token.AccessTTL))
	logger.Info("", zap.Duration("refresh token", cfg.Token.RefreshTTL))

	logger.Info("Server settings:")
	logger.Info("", zap.String("stage", cfg.App.Stage))
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
	if cfg.App.Settings.EnableAnalytics {
		logger.Info(fmt.Sprintf("Service: %s%s%s%s settings:", bold, line, "analytics", reset))
		logger.Info("", zap.String("address", cfg.Analytics.Host+":"+cfg.Analytics.Port))
	}
}
