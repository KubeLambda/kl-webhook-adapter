package app

import (
	"context"
	"go.uber.org/zap"
)

type loggerContextKey struct{}

func ContextWithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func Logger(ctx context.Context) *zap.SugaredLogger {
	logger := ctx.Value(loggerContextKey{}).(*zap.SugaredLogger)
	if logger == nil {
		panic("this context does not contain required SugaredLogger")
	}
	return logger
}

func BackgroundContextWithDefaultLogger() context.Context {
	return context.WithValue(context.Background(), loggerContextKey{}, zap.S())
}
