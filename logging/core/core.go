package core

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/feronatech/go-aws-logging/logging/core/utils"
)

const (
	LoggingContextKeyRegion      utils.LoggingContextKey = "region"
	LoggingContextKeyEnv         utils.LoggingContextKey = "env"
	LoggingContextKeyApplication utils.LoggingContextKey = "application"
	LoggingContextKeyLogLevel    utils.LoggingContextKey = "loglevel"
)

type LoggingOptions struct {
	Region      string
	Env         string
	Application string
	LogLevel    string
}

type Logger interface {
	SetContext(ctx context.Context) Logger
	Debug(message string, extra map[string]any)
	DebugCtx(ctx context.Context, message string, extra map[string]any)
	Info(message string, extra map[string]any)
	InfoCtx(ctx context.Context, message string, extra map[string]any)
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

func (l *logger) initializeWithWriter(writer io.Writer, logLevel string) Logger {
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     utils.MapLogLevelToSlogLevel(logLevel),
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: a.Key, Value: slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.000Z"))}
			}
			if a.Key == slog.MessageKey {
				return slog.Attr{Key: "message", Value: a.Value}
			}
			return a
		},
	})
	attrs := []slog.Attr{
		slog.String("region", l.region),
		slog.String("environment", l.env),
		slog.String("service", l.application),
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

func FromContext(ctx context.Context) Logger {
	l := &logger{
		region:      utils.GetValueFromContextAsStringOrDefault(ctx, LoggingContextKeyRegion, ""),
		env:         utils.GetValueFromContextAsStringOrDefault(ctx, LoggingContextKeyEnv, ""),
		application: utils.GetValueFromContextAsStringOrDefault(ctx, LoggingContextKeyApplication, ""),
		context:     ctx,
	}
	return l.initialize(utils.GetValueFromContextAsStringOrDefault(ctx, LoggingContextKeyLogLevel, "DEBUG"))
}

func (l *logger) SetContext(ctx context.Context) Logger {
	l.context = ctx
	return l
}

func handleExtraAttributes(extra map[string]any) slog.Attr {
	attrs := make([]slog.Attr, 0, len(extra))
	for key, value := range extra {
		attrs = append(attrs, slog.Any(key, value))
	}
	return slog.Any("context", attrs)
}

func (l *logger) Debug(message string, extra map[string]any) {
	l.InfoCtx(l.context, message, extra)
}

func (l *logger) DebugCtx(ctx context.Context, message string, extra map[string]any) {
	l.handler.LogAttrs(ctx, slog.LevelDebug, message, handleExtraAttributes(extra))
}

func (l *logger) Info(message string, extra map[string]any) {
	l.InfoCtx(l.context, message, extra)
}

func (l *logger) InfoCtx(ctx context.Context, message string, extra map[string]any) {
	l.handler.LogAttrs(ctx, slog.LevelInfo, message, handleExtraAttributes(extra))
}
