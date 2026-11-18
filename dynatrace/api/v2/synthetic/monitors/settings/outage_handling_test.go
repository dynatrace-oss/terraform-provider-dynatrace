//go:build unit

/**
* @license
* Copyright 2026 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package monitors

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const thresholdKey = "outage_handling.0.global_consecutive_outage_count_threshold"

// planThresholdDiff plans the given config against a state where the API returned the given threshold.
func planThresholdDiff(t *testing.T, stateThreshold string, config map[string]any) *terraform.InstanceDiff {
	t.Helper()
	resource := &schema.Resource{Schema: map[string]*schema.Schema{
		"outage_handling": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem:     &schema.Resource{Schema: new(OutageHandling).Schema()},
		},
	}}
	state := &terraform.InstanceState{
		ID: "id",
		Attributes: map[string]string{
			"id":                               "id",
			"outage_handling.#":                "1",
			"outage_handling.0.global_outages": "true",
			"outage_handling.0.local_outages":  "false",
			thresholdKey:                       stateThreshold,
		},
	}
	diff, err := resource.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]any{
		"outage_handling": []any{config},
	}), nil)
	require.NoError(t, err)
	return diff
}

func TestOutageHandling_GlobalThreshold_UnsetInConfig_DefaultInState_NoDiff(t *testing.T) {
	// HandlePreconditions sends 1 when the threshold is unset, so the refresh reads 1 back.
	diff := planThresholdDiff(t, "1", map[string]any{"global_outages": true})

	if diff != nil {
		assert.NotContains(t, diff.Attributes, thresholdKey)
	}
}

func TestOutageHandling_GlobalThreshold_UnsetInConfig_NonDefaultInState_Diff(t *testing.T) {
	// A value changed outside of Terraform must still be reverted to the default.
	diff := planThresholdDiff(t, "3", map[string]any{"global_outages": true})

	require.NotNil(t, diff)
	assert.Contains(t, diff.Attributes, thresholdKey)
}

func TestOutageHandling_GlobalThreshold_ChangedInConfig_Diff(t *testing.T) {
	diff := planThresholdDiff(t, "1", map[string]any{"global_outages": true, "global_consecutive_outage_count_threshold": 3})

	require.NotNil(t, diff)
	require.Contains(t, diff.Attributes, thresholdKey)
	assert.Equal(t, "3", diff.Attributes[thresholdKey].New)
}

func TestOutageHandling_HandlePreconditions_GlobalOutagesWithoutThreshold_SetsDefault(t *testing.T) {
	oh := OutageHandling{GlobalOutages: true}

	require.NoError(t, oh.HandlePreconditions())

	require.NotNil(t, oh.GlobalConsecutiveOutageCountThreshold)
	assert.Equal(t, defaultGlobalConsecutiveOutageCountThreshold, *oh.GlobalConsecutiveOutageCountThreshold)
}

func TestOutageHandling_HandlePreconditions_ThresholdWithoutOutages_Error(t *testing.T) {
	cases := []struct {
		name string
		oh   OutageHandling
	}{
		{name: "global", oh: OutageHandling{GlobalConsecutiveOutageCountThreshold: new(3)}},
		{name: "local consecutive", oh: OutageHandling{GlobalOutages: true, LocalConsecutiveOutageCountThreshold: new(3)}},
		{name: "local location", oh: OutageHandling{GlobalOutages: true, LocalLocationOutageCountThreshold: new(3)}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, tc.oh.HandlePreconditions())
		})
	}
}
