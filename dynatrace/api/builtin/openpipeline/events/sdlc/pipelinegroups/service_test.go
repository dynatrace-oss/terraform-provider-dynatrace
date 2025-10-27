package pipelinegroups_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

func TestOpenPipelineEventsSDLCPipelineGroups(t *testing.T) {
	t.Skip("Pipeline groups are disabled")
	api.TestAcc(t)
}
