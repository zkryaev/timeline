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
	"timeline/internal/repository/mail/notify"
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

func (a *App) Stop() {
	a.log.Info("Shutdown application...")
	a.httpServer.Shutdown(context.Background())
}

func (a *App) SetupControllers(tokenCfg config.Token, storage *postgres.PostgresRepo, mailService *notify.Mail /*redis*/) error {
	privateKey, err := secret.LoadPrivateKey()
	if err != nil {
		return err
	}
	usecaseAuth := auth.New(privateKey, storage, storage, tokenCfg, a.log)
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
	return nil
}
