/**
* @license
* Copyright 2025 Dynatrace LLC
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

package settings_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/routing/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var minimalRouting = settings.Routing{
	Kind: "events",
}

var pipelineId = "pipelineId"
var builtinPipelineId = "builtinPipelineId"
var routingWithEntries = settings.Routing{
	Kind: "events",
	RoutingEntries: []*settings.RoutingEntry{
		{
			Enabled:           true,
			PipelineType:      "custom",
			BuiltinPipelineID: nil,
			PipelineID:        &pipelineId,
			Matcher:           "some matcher",
			Description:       "somedescription",
		},
		{
			Enabled:           false,
			PipelineType:      "builtin",
			BuiltinPipelineID: &builtinPipelineId,
			PipelineID:        nil,
			Matcher:           "some matcher",
			Description:       "somedescription2",
		},
		{
			Enabled:           true,
			PipelineType:      "builtin",
			BuiltinPipelineID: &builtinPipelineId,
			PipelineID:        nil,
			Matcher:           "some other matcher",
			Description:       "somedescription3",
		},
	},
}

func TestRouting_MarshalHCL(t *testing.T) {
	cases := []struct {
		name     string
		input    settings.Routing
		expected hcl.Properties
	}{
		{
			name:  "minimum-set",
			input: minimalRouting,
			expected: hcl.Properties{
				"kind": "events",
			},
		},
		{
			name:  "with entries",
			input: routingWithEntries,
			expected: hcl.Properties{
				"kind": "events",
				"routing_entry": []interface{}{
					hcl.Properties{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "some matcher",
						"description":   "somedescription",
					},
					hcl.Properties{
						"enabled":             false,
						"builtin_pipeline_id": builtinPipelineId,
						"pipeline_type":       "builtin",
						"matcher":             "some matcher",
						"description":         "somedescription2",
					},
					hcl.Properties{
						"enabled":             true,
						"builtin_pipeline_id": builtinPipelineId,
						"pipeline_type":       "builtin",
						"matcher":             "some other matcher",
						"description":         "somedescription3",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var actual = hcl.Properties{}
			err := tc.input.MarshalHCL(actual)
			assert.Equal(t, tc.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestRouting_UnmarshalHCL(t *testing.T) {
	s := new(settings.Routing).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected settings.Routing
	}{
		{
			name: "minimal fields set",
			input: map[string]interface{}{
				"kind": "events",
			},
			expected: minimalRouting,
		},
		{
			name: "with entries",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "some matcher",
						"description":   "somedescription",
					},
					map[string]interface{}{
						"enabled":             false,
						"builtin_pipeline_id": builtinPipelineId,
						"pipeline_type":       "builtin",
						"matcher":             "some matcher",
						"description":         "somedescription2",
					},
					map[string]interface{}{
						"enabled":             true,
						"builtin_pipeline_id": builtinPipelineId,
						"pipeline_type":       "builtin",
						"matcher":             "some other matcher",
						"description":         "somedescription3",
					},
				},
			},
			expected: routingWithEntries,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual settings.Routing
			decoder := hcl.DecoderFrom(d)
			err := actual.UnmarshalHCL(decoder)
			assert.Equal(t, c.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestIngestSource_MarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    settings.Routing
		expected []byte
	}{
		{
			name:  "minimal-fields-set",
			input: minimalRouting,
			expected: []byte(
				`{
				}`),
		},
		{
			name:  "with entries",
			input: routingWithEntries,
			expected: []byte(
				`{
					"routingEntries": [
						{
							"enabled": true,	
							"pipelineId": "pipelineId",
							"pipelineType": "custom",
							"matcher": "some matcher",
							"description": "somedescription"
						},
						{
							"enabled": false,
							"builtinPipelineId": "builtinPipelineId",
							"pipelineType": "builtin",
							"matcher": "some matcher",
							"description": "somedescription2"
						},
						{
							"enabled": true,
							"builtinPipelineId": "builtinPipelineId",
							"pipelineType": "builtin",
							"matcher": "some other matcher",
							"description": "somedescription3"
						}
					]
				}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.input)
			require.NoError(t, err)

			var actualJSON map[string]interface{}
			err = json.Unmarshal(actual, &actualJSON)
			require.NoError(t, err)

			var expectedJSON map[string]interface{}
			err = json.Unmarshal(tc.expected, &expectedJSON)
			require.NoError(t, err)

			assert.Equal(t, expectedJSON, actualJSON)
		})
	}
}

func TestRouting_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			input: map[string]interface{}{
				"kind": "events",
			},
			wantErr: false,
		},
		{
			name: "valid with entries",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "some matcher",
						"description":   "somedescription",
					},
					map[string]interface{}{
						"enabled":             false,
						"builtin_pipeline_id": builtinPipelineId,
						"pipeline_type":       "builtin",
						"matcher":             "some matcher",
						"description":         "somedescription2",
					},
					map[string]interface{}{
						"enabled":             true,
						"builtin_pipeline_id": builtinPipelineId,
						"pipeline_type":       "builtin",
						"matcher":             "some matcher",
						"description":         "somedescription3",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "kind missing",
			input:   map[string]interface{}{},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "kind invalid",
			input: map[string]interface{}{
				"kind": "invalid",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf("expected kind to be one of %q, got invalid", settings.AllowedKinds),
		},
		{
			name: "description missing",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "some matcher",
					},
				},
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "description present but blank",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "some matcher",
						"description":   "",
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of routing_entry.0.description to be in the range (1 - 512), got ",
		},
		{
			name: "description starts with dt.",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "some matcher",
						"description":   "dt.asdf",
					},
				},
			},
			wantErr: true,
			errMsg:  "routing_entry.0.description must not start with 'dt.' or 'dynatrace.'",
		},
		{
			name: "description starts with dynatrace.",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "some matcher",
						"description":   "dynatrace.asdf",
					},
				},
			},
			wantErr: true,
			errMsg:  "routing_entry.0.description must not start with 'dt.' or 'dynatrace.'",
		},
		{
			name: "pipeline_type missing",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":     true,
						"pipeline_id": pipelineId,
						"matcher":     "some matcher",
						"description": "somedescription",
					},
				},
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "pipeline_type invalid",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "invalid",
						"matcher":       "some matcher",
						"description":   "somedescription",
					},
				},
			},
			wantErr: true,
			errMsg:  "expected routing_entry.0.pipeline_type to be one of [\"custom\" \"builtin\"], got invalid",
		},
		{
			name: "builtin_pipeline_id too long",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":             true,
						"builtin_pipeline_id": strings.Repeat("a", settings.BuiltinPipelineIDMaxLength+1),
						"pipeline_type":       "builtin",
						"matcher":             "some matcher",
						"description":         "somedescription",
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of routing_entry.0.builtin_pipeline_id to be in the range (1 - 500), got " + strings.Repeat("a", settings.BuiltinPipelineIDMaxLength+1),
		},
		{
			name: "matcher missing",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"description":   "somedescription",
					},
				},
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "matcher present but blank",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       "",
						"description":   "somedescription",
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of routing_entry.0.matcher to be in the range (1 - 1500), got ",
		},
		{
			name: "matcher too long",
			input: map[string]interface{}{
				"kind": "events",
				"routing_entry": []interface{}{
					map[string]interface{}{
						"enabled":       true,
						"pipeline_id":   pipelineId,
						"pipeline_type": "custom",
						"matcher":       strings.Repeat("a", settings.MatcherMaxLength+1),
						"description":   "somedescription",
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of routing_entry.0.matcher to be in the range (1 - 1500), got " + strings.Repeat("a", settings.MatcherMaxLength+1),
		},
	}

	r := &schema.Resource{Schema: new(settings.Routing).Schema()}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := terraform.NewResourceConfigRaw(tt.input)

			diags := r.Validate(c)
			assert.Equal(t, tt.wantErr, diags.HasError())

			if diags.HasError() {
				assert.Equal(t, tt.errMsg, diags[0].Summary)
			}
		})
	}
}
