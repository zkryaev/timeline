package app

import (
	"context"
	"fmt"
	"net/http"
	"timeline/internal/config"
	"timeline/internal/controller"
	authctrl "timeline/internal/controller/auth"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/domens/orgs"
	"timeline/internal/controller/domens/records"
	"timeline/internal/controller/domens/users"
	s3ctrl "timeline/internal/controller/s3"
	validation "timeline/internal/controller/validation"
	"timeline/internal/infrastructure"
	"timeline/internal/libs/secret"
	auth "timeline/internal/usecase/auth"
	"timeline/internal/usecase/orgcase"
	"timeline/internal/usecase/recordcase"
	s3usecase "timeline/internal/usecase/s3"
	"timeline/internal/usecase/usercase"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type App struct {
	httpServer http.Server
	log        *zap.Logger
}

func New(cfgApp config.Application, logger *zap.Logger) *App {
	app := &App{
		httpServer: http.Server{
			Addr:         cfgApp.Host + ":" + cfgApp.Port,
			ReadTimeout:  cfgApp.Timeout,
			WriteTimeout: cfgApp.Timeout,
			IdleTimeout:  cfgApp.IdleTimeout,
		},
		log: logger,
	}
	return app
}

func (a *App) Run() error {
	if err := a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to run server, %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) {
	a.httpServer.Shutdown(ctx)
}

func (a *App) SetupControllers(tokenCfg config.Token, storage infrastructure.Database, mailService infrastructure.Mail, s3Service infrastructure.S3 /*redis*/) error {
	privateKey, err := secret.LoadPrivateKey()
	if err != nil {
		return err
	}
	// Инициализация Auth
	usecaseAuth := auth.New(
		privateKey,
		storage,
		storage,
		storage,
		mailService,
		tokenCfg,
		a.log,
	)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	validator := validation.NewCustomValidator()

	authAPI := authctrl.New(
		usecaseAuth,
		middleware.New(privateKey, a.log),
		a.log,
		json,
		validator,
	)

	// Инициализация User
	usecaseUser := usercase.New(
		storage,
		storage,
		storage,
		a.log,
	)

	userAPI := users.NewUserCtrl(
		usecaseUser,
		a.log,
		json,
		validator,
	)

	// Инициализация Org
	usecaseOrg := orgcase.New(
		storage,
		storage,
		a.log,
	)

	orgAPI := orgs.NewOrgCtrl(
		usecaseOrg,
		a.log,
		json,
		validator,
	)

	// Работа с записями
	usecaseRecord := recordcase.New(
		storage,
		storage,
		storage,
		mailService,
		a.log,
	)

	recordAPI := records.NewRecordCtrl(
		usecaseRecord,
		a.log,
		json,
		validator,
	)

	// Хранилище изображений
	s3Usecase := s3usecase.New(storage, storage, s3Service, a.log)
	s3API := s3ctrl.New(s3Usecase, a.log, json)
	controllerSet := &controller.Controllers{
		Auth:   authAPI,
		User:   userAPI,
		Org:    orgAPI,
		Record: recordAPI,
		S3:     s3API,
	}

	a.httpServer.Handler = controller.InitRouter(controllerSet)
	return nil
}
