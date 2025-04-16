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
	"timeline/internal/controller/metrics"
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
	logger.Info("Application settings:")
	logger.Info("", zap.String("environment", cfg.App.Env))
	logger.Info("", zap.Bool("use_local_backdata", cfg.App.Settings.UseLocalBackData))
	logger.Info("", zap.Bool("enable_authorization", cfg.App.Settings.EnableAuthorization))
	logger.Info("", zap.Bool("enable_repo_s3", cfg.App.Settings.EnableRepoS3))
	logger.Info("", zap.Bool("enable_repo_mail", cfg.App.Settings.EnableRepoMail))
	logger.Info("", zap.Bool("enable_metrics", cfg.App.Settings.EnableMetrics))
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
	logger.Info(
		"DB",
		zap.String("database server", cfg.DB.Host+":"+cfg.DB.Port),
		zap.String("ssl", cfg.DB.SSLmode),
	)
	defer db.Close()

	backdata := &loader.BackData{}
	if cfg.App.Settings.UseLocalBackData {
		logger.Info("Loading from local storage")
		logger.Info("Start loading from DB")
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
	if cfg.App.Settings.EnableRepoMail {
		// Поднимаем почтовый сервис параметрами по умолчанию
		post = mail.New(cfg.Mail, logger, 0, 0, 0)
		post.Start()
		logger.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("MAIL_HOST")))
		logger.Info(
			"MAIL",
			zap.String("mail server", cfg.Mail.Host+":"+strconv.Itoa(cfg.Mail.Port)),
		)
		defer post.Shutdown()
	} else {
		logger.Info("Mail launch skipped")
	}

	var s3repo *s3.Minio
	if cfg.App.Settings.EnableRepoMail {
		// Подключение к S3
		s3repo = s3.New(cfg.S3)
		if err = s3repo.Connect(); err != nil {
			logger.Fatal(fmt.Sprintf("failed to connect to %s", os.Getenv("S3")), zap.Error(err))
		}
		logger.Info(fmt.Sprintf("%s %s", successConnection, os.Getenv("S3")))
		logger.Info(
			"S3",
			zap.String("storage", cfg.S3.Host+":"+cfg.S3.DataPort),
			zap.String("console", cfg.S3.Host+":"+cfg.S3.ConsolePort),
			zap.Bool("ssl", cfg.S3.SSLmode),
		)
	} else {
		logger.Info("S3 launch skipped")
	}

	s := cronjob.InitCronScheduler(db)
	defer s.Shutdown()
	s.Start()

	ctx, cancel := context.WithCancel(context.Background())
	timeout := 1 * time.Minute
	errch := make(chan error, 2)
	app := app.New(cfg.App, logger)
	err = app.SetupControllers(cfg.Token, backdata, db, post, s3repo)
	if err != nil {
		logger.Fatal(
			"failed setup controllers",
			zap.Error(err),
		)
	}
	app.Run(errch)
	defer app.Shutdown(ctx, timeout)
	logger.Info("Application is running")
	logger.Info("", zap.String("listening on", cfg.App.Server.Host+":"+cfg.App.Server.Port))

	if cfg.App.Settings.EnableMetrics {
		metricsExporter := metrics.NewPrometheusExporter(cfg.Prometheus, logger)
		metricsExporter.Launch(errch)
		defer metricsExporter.Shutdown(ctx, timeout)
		logger.Info("Prometheus exporter launched")
		logger.Info("", zap.String("listening on", cfg.Prometheus.Host+":"+cfg.Prometheus.Port))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		cancel()
		logger.Info("Received signal", zap.String("signal", sig.String()))
	case err = <-errch:
		cancel()
		logger.Error("error occurred", zap.Error(err))
	}
	logger.Info("Application stopped")
}
