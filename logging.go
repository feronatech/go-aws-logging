package go_aws_logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

type LoggingOptions struct {
	Region      string
	Env         string
	Application string
	LogLevel    string
}

type Logger interface {
	SetContext(ctx context.Context) Logger
	Info(message string)
	InfoCtx(ctx context.Context, message string)
}

type logger struct {
	region      string
	env         string
	application string
	context     context.Context
	handler     *slog.Logger
}

func (l *logger) initialize(logLevel string) Logger {
	return l.initializeWithWriter(os.Stdout, logLevel)
}

func mapLogLevelToSlogLevel(logLevel string) slog.Level {
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	}
	return 0 // defaults to INFO
}

func (l *logger) initializeWithWriter(writer io.Writer, logLevel string) Logger {
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     mapLogLevelToSlogLevel(logLevel),
		AddSource: false,
	})
	attrs := []slog.Attr{
		slog.String("region", l.region),
		slog.String("env", l.env),
		slog.String("application", l.application),
	}
	l.handler = slog.New(handler.WithAttrs(attrs))
	return l
}

func NewLogger(region string, env string, application string) Logger {
	l := &logger{
		region:      region,
		env:         env,
		application: application,
		context:     context.Background(),
	}
	return l.initialize(slog.LevelDebug.String())
}

func FromOptions(options LoggingOptions) Logger {
	l := &logger{
		region:      options.Region,
		env:         options.Env,
		application: options.Application,
		context:     context.Background(),
	}
	return l.initialize(options.LogLevel)
}

func getValueFromContextAsStringOrDefault(ctx context.Context, key string, defaultValue string) string {
	if ctx == nil {
		return defaultValue
	}
	val := ctx.Value(key)
	if strVal, ok := val.(string); ok {
		return strVal
	}
	return defaultValue
}

func FromContext(ctx context.Context) Logger {
	l := &logger{
		region:      getValueFromContextAsStringOrDefault(ctx, "region", ""),
		env:         getValueFromContextAsStringOrDefault(ctx, "env", ""),
		application: getValueFromContextAsStringOrDefault(ctx, "application", ""),
		context:     ctx,
	}
	return l.initialize(getValueFromContextAsStringOrDefault(ctx, "loglevel", "DEBUG"))
}

func (l *logger) SetContext(ctx context.Context) Logger {
	l.context = ctx
	return l
}

func (l *logger) Info(message string) {
	l.InfoCtx(l.context, message)
}

func (l *logger) InfoCtx(ctx context.Context, message string) {
	l.handler.InfoContext(ctx, message)
}
