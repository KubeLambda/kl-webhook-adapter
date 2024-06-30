package infra

import (
	"context"
	"go.uber.org/zap"
	"serverless-service-webhook-adapter/internal/adapters/apiserver"
	"serverless-service-webhook-adapter/internal/core/app"
)

func Start(deployment string) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	ctx := app.ContextWithLogger(context.Background(), zap.S())

	cfg := app.LoadConfig(deployment)
	apiserver.Start(ctx, cfg)
}
