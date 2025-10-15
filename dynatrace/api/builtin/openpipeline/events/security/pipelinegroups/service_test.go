package pipelinegroups_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

func TestOpenPipelineEventsSecurityPipelineGroups(t *testing.T) {
	t.Skip("Feature not enabled")

	api.TestAcc(t)
}
