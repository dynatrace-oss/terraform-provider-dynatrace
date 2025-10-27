package pipelinegroups_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

func TestOpenPipelineUserEventsPipelineGroups(t *testing.T) {
	t.Skip("Pipeline groups are disabled")
	api.TestAcc(t)
}
