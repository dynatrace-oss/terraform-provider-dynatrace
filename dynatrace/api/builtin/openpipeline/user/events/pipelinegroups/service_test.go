package pipelinegroups_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	featureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/provider/featureflag"
)

func TestOpenPipelineUserEventsPipelineGroups(t *testing.T) {
	if !featureflags.OpenPipelinePipelineGroups.Enabled() {
		t.Skip("Pipeline groups are disabled")
	}
	api.TestAcc(t)
}
