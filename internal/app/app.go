package app

import (
	"context"
	"fmt"
	"net/http"
	"timeline/internal/config"
	"timeline/internal/controller"
	authctrl "timeline/internal/controller/auth"
	"timeline/internal/libs/secret"
	"timeline/internal/repository/database/postgres"
	auth "timeline/internal/usecase/auth"
	"timeline/internal/usecase/auth/middleware"

	"github.com/go-playground/validator"
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
			Addr:         cfgApp.Host + cfgApp.Port,
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
		a.log.Fatal("failed to run server")
		return fmt.Errorf("failed to run server, %w", err)
	}
	return nil
}

func (a *App) Stop() {
	a.log.Info("Shutdown application...")
	a.httpServer.Shutdown(context.Background())
	a.log.Info("App stopped")
}

func (a *App) SetupControllers(tokenCfg config.Token, db *postgres.PostgresRepo /*redis*/) {
	// usecase ->
	// controller
	// TODO: добавить логирование
	privateKey, err := secret.LoadPrivateKey("TODO")
	if err != nil {
		panic(err) // TODO: а как иначе елки палки
	}
	usecaseAuth := auth.New(privateKey, db, tokenCfg, a.log)
	authAPI := authctrl.New(
		usecaseAuth,
		middleware.New(privateKey),
		a.log,
		jsoniter.ConfigCompatibleWithStandardLibrary,
		*validator.New(),
	)

	controllerSet := &controller.Controllers{
		Auth: authAPI,
	}

	a.httpServer.Handler = controller.InitRouter(controllerSet)
}
