package go_aws_logging

import (
	"context"
	"log/slog"
)

type Logger struct {
	region      string
	env         string
	application string
}

func NewLogger(region, env, application string) *Logger {
	return &Logger{
		region:      region,
		env:         env,
		application: application,
	}
}

func (l *Logger) Info(message string) {
	l.InfoCtx(context.Background(), message)
}

func (l *Logger) InfoCtx(ctx context.Context, message string) {
	slog.InfoContext(ctx, message, "region", l.region, "env", l.env, "application", l.application)
}
