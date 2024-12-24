package app

import (
	"context"
	"fmt"
	"net/http"
	"timeline/internal/config"
	"timeline/internal/controller"
	authctrl "timeline/internal/controller/auth"
	"timeline/internal/controller/domens/orgs"
	"timeline/internal/controller/domens/records"
	"timeline/internal/controller/domens/users"
	validation "timeline/internal/controller/validation"
	"timeline/internal/libs/secret"
	"timeline/internal/repository"
	"timeline/internal/repository/mail"
	auth "timeline/internal/usecase/auth"
	"timeline/internal/usecase/auth/middleware"
	"timeline/internal/usecase/orgcase"
	"timeline/internal/usecase/recordcase"
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

func (a *App) SetupControllers(tokenCfg config.Token, storage repository.Repository, mailService mail.Post /*redis*/) error {
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

	controllerSet := &controller.Controllers{
		Auth:   authAPI,
		User:   userAPI,
		Org:    orgAPI,
		Record: recordAPI,
	}

	a.httpServer.Handler = controller.InitRouter(controllerSet)
	return nil
}
