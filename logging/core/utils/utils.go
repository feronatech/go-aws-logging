package utils

import (
	"context"
	"log/slog"
	"strings"
)

type LoggingContextKey string

func MapLogLevelToSlogLevel(logLevel string) slog.Level {
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

func GetValueFromContextAsStringOrDefault(ctx context.Context, key LoggingContextKey, defaultValue string) string {
	if ctx == nil {
		return defaultValue
	}
	val := ctx.Value(key)
	if strVal, ok := val.(string); ok {
		return strVal
	}
	return defaultValue
}
