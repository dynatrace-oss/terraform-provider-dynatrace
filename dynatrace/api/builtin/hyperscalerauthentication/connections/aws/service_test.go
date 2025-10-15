package aws_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

func TestAwsConnection(t *testing.T) {
	t.Skip("Test is skipped as long as no AWS test environment is setup to do a proper E2E test")

	api.TestAcc(t)
}
