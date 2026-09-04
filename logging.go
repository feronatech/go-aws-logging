package go_aws_logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type LoggingOptions struct {
	Region      string
	Env         string
	Application string
}

type logger struct {
	region      string
	env         string
	application string
	context     context.Context
	handler     *slog.Logger
}

func (l *logger) initialize() *logger {
	return l.initializeWithWriter(os.Stdout)
}

func (l *logger) initializeWithWriter(writer io.Writer) *logger {
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
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

func NewLogger(region string, env string, application string) *logger {
	logger := &logger{
		region:      region,
		env:         env,
		application: application,
		context:     context.Background(),
	}
	return logger.initialize()
}

func FromOptions(options LoggingOptions) *logger {
	logger := &logger{
		region:      options.Region,
		env:         options.Env,
		application: options.Application,
		context:     context.Background(),
	}
	return logger.initialize()
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

func FromContext(ctx context.Context) *logger {
	logger := &logger{
		region:      getValueFromContextAsStringOrDefault(ctx, "region", ""),
		env:         getValueFromContextAsStringOrDefault(ctx, "env", ""),
		application: getValueFromContextAsStringOrDefault(ctx, "application", ""),
		context:     ctx,
	}
	return logger.initialize()
}

func FromConfig(config *AwsConfig) *logger {
	logger := &logger{
		region:      config.Region,
		env:         config.RuntimeEnvironment,
		application: config.AppID,
		context:     context.Background(),
	}
	return logger.initialize()
}

func (l *logger) WithContext(ctx context.Context) *logger {
	newLogger := &logger{
		region:      l.region,
		env:         l.env,
		application: l.application,
		context:     ctx,
	}
	return newLogger.initialize()
}

func (l *logger) Info(message string) {
	l.InfoCtx(l.context, message)
}

func (l *logger) InfoCtx(ctx context.Context, message string) {
	l.handler.InfoContext(ctx, message)
}
