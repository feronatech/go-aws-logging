package go_aws_logging

import (
	"bytes"
	"context"
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

func (l *logger) testLogger(t *testing.T) (*logger, LogGetter) {
	t.Helper()
	writer := &bytes.Buffer{}
	return l.initializeWithWriter(writer), &logGetter{buffer: writer}
}

func Test_LoggerInitialize(t *testing.T) {
	regionInput := "us-east-1"
	envInput := "dev"
	applicationInput := "my-app"
	logger := NewLogger(regionInput, envInput, applicationInput)

	if logger.region != regionInput {
		t.Errorf("Expected region to be '%s', got '%s'", regionInput, logger.region)
	}
	if logger.env != envInput {
		t.Errorf("Expected env to be '%s', got '%s'", envInput, logger.env)
	}
	if logger.application != applicationInput {
		t.Errorf("Expected application to be '%s', got '%s'", applicationInput, logger.application)
	}
	if logger.handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
	if logger.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}

func Test_LoggerFromOptions(t *testing.T) {
	options := LoggingOptions{
		Region:      "us-west-2",
		Env:         "prod",
		Application: "my-service",
	}
	logger := FromOptions(options)

	if logger.region != options.Region {
		t.Errorf("Expected region to be '%s', got '%s'", options.Region, logger.region)
	}
	if logger.env != options.Env {
		t.Errorf("Expected env to be '%s', got '%s'", options.Env, logger.env)
	}
	if logger.application != options.Application {
		t.Errorf("Expected application to be '%s', got '%s'", options.Application, logger.application)
	}
	if logger.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}

func Test_LoggerFromEmptyContext(t *testing.T) {
	ctx := context.Background()
	logger := FromContext(ctx)

	if logger.region != "" {
		t.Errorf("Expected region to be empty, got '%s'", logger.region)
	}
	if logger.env != "" {
		t.Errorf("Expected env to be empty, got '%s'", logger.env)
	}
	if logger.application != "" {
		t.Errorf("Expected application to be empty, got '%s'", logger.application)
	}
	if logger.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}

func Test_LoggerFromContextWithValues(t *testing.T) {
	regionInput := "eu-central-1"
	envInput := "staging"
	applicationInput := "test-app"

	ctx := context.WithValue(context.Background(), "region", regionInput)
	ctx = context.WithValue(ctx, "env", envInput)
	ctx = context.WithValue(ctx, "application", applicationInput)
	logger := FromContext(ctx)

	if logger.region != regionInput {
		t.Errorf("Expected region to be '%s', got '%s'", regionInput, logger.region)
	}
	if logger.env != envInput {
		t.Errorf("Expected env to be '%s', got '%s'", envInput, logger.env)
	}
	if logger.application != applicationInput {
		t.Errorf("Expected application to be '%s', got '%s'", applicationInput, logger.application)
	}
	if logger.context == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}
