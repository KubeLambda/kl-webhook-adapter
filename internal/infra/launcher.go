package infra

import (
	"context"
	"kl-webhook-adapter/internal/adapters/apiserver"
	"kl-webhook-adapter/internal/core/app"

	"go.uber.org/zap"
)

func Start(deployment string) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	ctx := app.ContextWithLogger(context.Background(), zap.S())

	cfg := app.LoadConfig(deployment)
	apiserver.Start(ctx, cfg)
}
