package go_aws_logging

type AwsConfig struct {
	Region             string `readonly:"true"`
	DefaultsMode       string `readonly:"true"`
	RuntimeEnvironment string `readonly:"true"`
	AppID              string `readonly:"true"`
}
