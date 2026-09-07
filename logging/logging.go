package logging

import (
	"github.com/feronatech/go-aws-logging/logging/core"
)

type Logger = core.Logger
type LoggingOptions = core.LoggingOptions

var LoggingContextKeyRegion = core.LoggingContextKeyRegion
var LoggingContextKeyEnv = core.LoggingContextKeyEnv
var LoggingContextKeyApplication = core.LoggingContextKeyApplication
var LoggingContextKeyLogLevel = core.LoggingContextKeyLogLevel

var NewLogger = core.NewLogger
var FromOptions = core.FromOptions
var FromContext = core.FromContext
