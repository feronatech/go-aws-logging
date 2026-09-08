// logging package provides the entry interface to the library.
// if you import this package, you will usually get no peer dependencies in the final binary.
package logging

import (
	"github.com/feronatech/go-aws-logging/logging/core"
)

// Logger is the main interface.
// It provides methods for setting context and logging messages at different levels (Debug, Info, Warn, Error).
// This is what you will get when calling NewLogger, FromOptions, or FromContext.
// And it is typically what you will use in your application code to log messages.
type Logger = core.Logger

// LoggingOptions is a struct that holds configuration options for the logger.
// It includes fields for Region, Env, Application, and LogLevel.
// It is used to initialize a new logger with specific settings.
type LoggingOptions = core.LoggingOptions

// LoggingContextKeyRegion is a context key for the region.
// It is used to store and retrieve the region value from a context.Context.
var LoggingContextKeyRegion = core.LoggingContextKeyRegion

// LoggingContextKeyEnv is a context key for the environment.
// It is used to store and retrieve the environment value from a context.Context.
var LoggingContextKeyEnv = core.LoggingContextKeyEnv

// LoggingContextKeyApplication is a context key for the application name.
// It is used to store and retrieve the application name value from a context.Context.
var LoggingContextKeyApplication = core.LoggingContextKeyApplication

// LoggingContextKeyLogLevel is a context key for the log level.
// It is used to store and retrieve the log level value from a context.Context.
var LoggingContextKeyLogLevel = core.LoggingContextKeyLogLevel

// NewLogger creates a new logger with the specified region, environment, and application name.
// It initializes the logger with default settings and returns it as a Logger interface.
var NewLogger = core.NewLogger

// FromOptions creates a new logger based on the provided LoggingOptions.
// It initializes the logger with the specified settings and returns it as a Logger interface.
var FromOptions = core.FromOptions

// FromContext creates a new logger based on the values stored in the provided context.Context.
// It retrieves the region, environment, application name, and log level from the context and initializes the logger accordingly.
// If any of these values are not present in the context, they will default to empty strings.
var FromContext = core.FromContext
