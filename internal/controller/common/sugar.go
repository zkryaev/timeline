package common

import (
	"context"
	"timeline/internal/controller/scope"

	"go.uber.org/zap"
)

func LoggerWithUUID(settings *scope.Settings, Logger *zap.Logger, ctx context.Context) *zap.Logger {
	if !settings.EnableAuthorization {
		return Logger
	}
	uuid, _ := ctx.Value("uuid").(string)
	logger := Logger.With(zap.String("uuid", uuid))
	return logger
}
