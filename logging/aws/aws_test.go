package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func Test_LoggerFromConfig(t *testing.T) {
	config := &aws.Config{
		Region: "us-west-2",
		AppID:  "test-app",
	}
	l := FromConfig(config)

	if l == nil {
		t.Error("Expected context to be initialized, got nil")
	}
}
