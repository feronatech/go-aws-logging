package utils

import (
	"context"
	"log/slog"
	"testing"
)

func Test_MapLogLevelToSlogLevel(t *testing.T) {
	tests := []struct {
		logLevel string
		expected slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"WARN", slog.LevelWarn},
		{"ERROR", slog.LevelError},
		{"debug", slog.LevelDebug},
		{"unknown", 0}, // default case
	}
	for _, tt := range tests {
		t.Run(tt.logLevel, func(t *testing.T) {
			if got := MapLogLevelToSlogLevel(tt.logLevel); got != tt.expected {
				t.Errorf("MapLogLevelToSlogLevel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_GetValueFromContextAsStringOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		key          LoggingContextKey
		defaultValue string
		want         string
	}{
		{
			name:         "value present",
			ctx:          context.WithValue(context.Background(), LoggingContextKey("test"), "value"),
			key:          LoggingContextKey("test"),
			defaultValue: "default",
			want:         "value",
		},
		{
			name:         "value not present",
			ctx:          context.Background(),
			key:          LoggingContextKey("test"),
			defaultValue: "default",
			want:         "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValueFromContextAsStringOrDefault(tt.ctx, tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("GetValueFromContextAsStringOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
