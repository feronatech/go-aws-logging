// logging/aws package provides a logger implementation that is compatible with AWS SDK v2.
// This library must be explicitly imported with "github.com/feronatech/go-aws-logging/logging/aws".
// The design of this library was defined on purpose to be able to have the aws package as a peer dependency of the logging package
// ending up in the final binary only when this package is used.
package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/feronatech/go-aws-logging/logging"
)

// FromConfig creates a new logger based on the provided AWS SDK v2 configuration.
// It retrieves the region, environment, and application name from the configuration and initializes the logger accordingly.
// If any of these values are not present in the configuration, they will default to empty strings.
func FromConfig(config *aws.Config) logging.Logger {
	return logging.NewLogger(
		config.Region,
		string(config.RuntimeEnvironment.EnvironmentIdentifier),
		config.AppID,
	)
}
