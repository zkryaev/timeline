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
	"timeline/internal/controller/monitoring"
	s3ctrl "timeline/internal/controller/s3"
	"timeline/internal/controller/scope"
	validation "timeline/internal/controller/validation"
	"timeline/internal/infrastructure"
	"timeline/internal/sugar/secret"
	auth "timeline/internal/usecase/auth"
	"timeline/internal/usecase/orgcase"
	"timeline/internal/usecase/recordcase"
	s3usecase "timeline/internal/usecase/s3"
	"timeline/internal/usecase/usercase"
	"timeline/internal/utils/loader"

	"go.uber.org/zap"
)

type App struct {
	httpServer *http.Server
	log        *zap.Logger
	appcfg     config.Application
}

func New(cfgApp config.Application, logger *zap.Logger) *App {
	return &App{
		appcfg: cfgApp,
		httpServer: &http.Server{
			Addr:         cfgApp.Host + ":" + cfgApp.Port,
			ReadTimeout:  cfgApp.Timeout,
			WriteTimeout: cfgApp.Timeout,
			IdleTimeout:  cfgApp.IdleTimeout,
		},
		log: logger,
	}
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

func (a *App) SetupControllers(tokenCfg config.Token, backdata *loader.BackData, storage infrastructure.Database, mailService infrastructure.Mail, s3Service infrastructure.S3) error {
	privateKey, err := secret.LoadPrivateKey()
	if err != nil {
		return err
	}

	validator, err := validation.NewCustomValidator() // singleton is right way to use it (thread-safe!)
	if err != nil {
		return err
	}

	settings := scope.NewDefaultSettings(a.appcfg)
	routes := scope.NewDefaultRoutes(settings)
	middleware := middleware.New(privateKey, a.log, routes)

	monitorAPI := monitoring.New(a.log, settings)

	authAPI := authctrl.New(
		auth.New(
			privateKey,
			storage,
			storage,
			storage,
			mailService,
			tokenCfg,
			settings,
		),
		middleware,
		a.log,
		settings,
	)
	var s3API *s3ctrl.S3Ctrl
	if settings.EnableMedia {
		s3API = s3ctrl.New(
			s3usecase.New(
				storage,
				storage,
				s3Service,
			),
			a.log,
			settings,
		)
	}

	userAPI := users.New(
		usercase.New(
			storage,
			storage,
			storage,
			backdata,
			settings,
		),
		a.log,
		validator,
		middleware,
		settings,
	)

	orgAPI := orgs.New(
		orgcase.New(
			storage,
			storage,
			backdata,
		),
		middleware,
		a.log,
		settings,
	)

	recordAPI := records.New(
		recordcase.New(
			backdata,
			storage,
			storage,
			storage,
			mailService,
			settings,
		),
		middleware,
		a.log,
		settings,
	)

	controllerSet := &controller.Controllers{
		Monitor: monitorAPI,
		Auth:    authAPI,
		User:    userAPI,
		Org:     orgAPI,
		Record:  recordAPI,
		S3:      s3API,
	}
	monitorAPI.Router = controller.InitRouter(controllerSet, routes, settings)
	a.httpServer.Handler = monitorAPI.Router
	return nil
}
