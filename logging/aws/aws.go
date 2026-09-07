package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/feronatech/go-aws-logging/logging"
)

func FromConfig(config *aws.Config) logging.Logger {
	return logging.NewLogger(
		config.Region,
		string(config.RuntimeEnvironment.EnvironmentIdentifier),
		config.AppID,
	)
}
