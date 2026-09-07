package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

type LogGetter interface {
	Get() string
}

type logGetter struct {
	buffer *bytes.Buffer
}

func (l *logGetter) Get() string {
	return l.buffer.String()
}

func (l *logger) testLogger(t *testing.T) (logger, LogGetter) {
	t.Helper()
	writer := &bytes.Buffer{}
	if log, ok := l.initializeWithWriter(writer, slog.LevelDebug.String()).(*logger); ok {
		return *log, &logGetter{buffer: writer}
	}
	t.Fail()
	return logger{}, nil
}

func Test_LoggerInitialize(t *testing.T) {
	regionInput := "us-east-1"
	envInput := "dev"
	applicationInput := "my-app"
	l := NewLogger(regionInput, envInput, applicationInput).(*logger)

	if l.region != regionInput {
		t.Errorf("Expected region to be '%s', got '%s'", regionInput, l.region)
	}
	if l.env != envInput {
		t.Errorf("Expected env to be '%s', got '%s'", envInput, l.env)
	}
	if l.application != applicationInput {
		t.Errorf("Expected application to be '%s', got '%s'", applicationInput, l.application)
	}
	if l.handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
	if l.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}

func Test_LoggerFromOptions(t *testing.T) {
	options := LoggingOptions{
		Region:      "us-west-2",
		Env:         "prod",
		Application: "my-service",
	}
	l := FromOptions(options).(*logger)

	if l.region != options.Region {
		t.Errorf("Expected region to be '%s', got '%s'", options.Region, l.region)
	}
	if l.env != options.Env {
		t.Errorf("Expected env to be '%s', got '%s'", options.Env, l.env)
	}
	if l.application != options.Application {
		t.Errorf("Expected application to be '%s', got '%s'", options.Application, l.application)
	}
	if l.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}

func Test_LoggerFromEmptyContext(t *testing.T) {
	ctx := context.Background()
	l := FromContext(ctx).(*logger)

	if l.region != "" {
		t.Errorf("Expected region to be empty, got '%s'", l.region)
	}
	if l.env != "" {
		t.Errorf("Expected env to be empty, got '%s'", l.env)
	}
	if l.application != "" {
		t.Errorf("Expected application to be empty, got '%s'", l.application)
	}
	if l.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}

func Test_LoggerFromContextWithValues(t *testing.T) {
	regionInput := "eu-central-1"
	envInput := "staging"
	applicationInput := "test-app"

	ctx := context.WithValue(context.Background(), LoggingContextKeyRegion, regionInput)
	ctx = context.WithValue(ctx, LoggingContextKeyEnv, envInput)
	ctx = context.WithValue(ctx, LoggingContextKeyApplication, applicationInput)
	l := FromContext(ctx).(*logger)

	if l.region != regionInput {
		t.Errorf("Expected region to be '%s', got '%s'", regionInput, l.region)
	}
	if l.env != envInput {
		t.Errorf("Expected env to be '%s', got '%s'", envInput, l.env)
	}
	if l.application != applicationInput {
		t.Errorf("Expected application to be '%s', got '%s'", applicationInput, l.application)
	}
	if l.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}

func Test_LoggerSetContext(t *testing.T) {
	regionInput := "ap-southeast-1"
	envInput := "qa"
	applicationInput := "demo-app"
	regionValue := "new-region"
	l := NewLogger(regionInput, envInput, applicationInput).(*logger)
	if l.context.Value(LoggingContextKeyRegion) != nil {
		t.Error("Expected context to be empty, got a value")
	}
	ctx := context.WithValue(context.Background(), LoggingContextKeyRegion, regionValue)
	l = l.SetContext(ctx).(*logger)
	if l.context == nil {
		t.Error("Expected context to be set, got nil")
	}
	if l.context.Value(LoggingContextKeyRegion) != regionValue {
		t.Errorf("Expected context region to be '%s', got '%v'", regionValue, l.context.Value(LoggingContextKeyRegion))
	}
}

func Test_LoggerInfo(t *testing.T) {
	regionInput := "us-east-1"
	envInput := "dev"
	applicationInput := "my-app"
	l := NewLogger(regionInput, envInput, applicationInput).(*logger)
	testLogger, logGetter := l.testLogger(t)

	testMessage := "This is a test log message"
	testLogger.Info(testMessage)

	logOutput := logGetter.Get()
	if !bytes.Contains([]byte(logOutput), []byte(testMessage)) {
		t.Errorf("Expected log output to contain '%s', got '%s'", testMessage, logOutput)
	}
}
