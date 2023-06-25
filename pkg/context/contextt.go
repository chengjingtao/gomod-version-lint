package context

import (
	"context"
	"go.uber.org/zap"
)

var loggerKey = struct{}{}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.SugaredLogger {
	return ctx.Value(loggerKey).(*zap.SugaredLogger)
}
