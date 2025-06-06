package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"timeline/internal/config"
	"timeline/internal/controller"
	authctrl "timeline/internal/controller/auth"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/domens/orgs"
	"timeline/internal/controller/domens/records"
	"timeline/internal/controller/domens/users"
	"timeline/internal/controller/external"
	"timeline/internal/controller/monitoring"
	"timeline/internal/controller/monitoring/metrics"
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
	server     *http.Server
	log        *zap.Logger
	appcfg     config.Application
	serverOnce sync.Once
	wg         sync.WaitGroup
}

func New(cfgApp config.Application, logger *zap.Logger) *App {
	return &App{
		appcfg: cfgApp,
		server: &http.Server{
			Addr:         cfgApp.Server.Host + ":" + cfgApp.Server.Port,
			ReadTimeout:  cfgApp.Server.Timeout,
			WriteTimeout: cfgApp.Server.Timeout,
			IdleTimeout:  cfgApp.Server.IdleTimeout,
		},
		log: logger,
	}
}

func (a *App) SetHandler(handler http.Handler) {
	a.server.Handler = handler
}

func (a *App) Run(errch chan error) {
	a.serverOnce.Do(func() {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errch <- fmt.Errorf("backend server error: %w", err)
			}
		}()
	})
}

func (a *App) Shutdown(cancelCtx context.Context, timeout time.Duration) {
	timeoutCtx, cancel := context.WithTimeout(cancelCtx, timeout)
	defer cancel()
	if err := a.server.Shutdown(timeoutCtx); err != nil {
		a.log.Error("failed to shutdown HTTP server", zap.Error(err))
	}
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		a.log.Info("Backend server closing complete")
	case <-timeoutCtx.Done():
		a.log.Error("timeout while closing HTTP server", zap.Error(timeoutCtx.Err()))
	}
}

func (a *App) SetupControllers(cfg *config.Config, backdata *loader.BackData, storage infrastructure.Database, mailService infrastructure.Mail, s3Service infrastructure.S3) error {
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
	var metricList *metrics.Metrics
	if settings.EnableMetrics {
		metricList = metrics.NewDefaultMetrics()
	}
	middleware := middleware.New(privateKey, a.log, routes, metricList)

	monitorAPI := monitoring.New(a.log, settings)

	authAPI := authctrl.New(
		auth.New(
			privateKey,
			storage,
			storage,
			storage,
			mailService,
			cfg.Token,
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
		settings,
	)

	orgAPI := orgs.New(
		orgcase.New(
			storage,
			storage,
			backdata,
		),
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
		a.log,
		settings,
	)

	analyticsRedirector := external.NewAnalyticsClient(cfg.Analytics, a.log, settings)

	controllerSet := &controller.Controllers{
		Monitor:   monitorAPI,
		Auth:      authAPI,
		User:      userAPI,
		Org:       orgAPI,
		Record:    recordAPI,
		S3:        s3API,
		Analitycs: analyticsRedirector,
	}
	monitorAPI.Router = controller.InitRouter(controllerSet, routes, settings)
	a.SetHandler(monitorAPI.Router)
	return nil
}
