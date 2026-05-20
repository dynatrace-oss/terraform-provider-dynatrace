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

package metricevents

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestQueryDefinition_HandlePreconditions_Valid(t *testing.T) {
	cases := []struct {
		name  string
		input map[string]any
	}{
		{
			name: "MetricKey_Minimal",
			input: map[string]any{
				"type":        "METRIC_KEY",
				"aggregation": "AVG",
				"metric_key":  "builtin:host.cpu.usage",
			},
		},
		{
			name: "MetricKey_WithDimensionFilter",
			input: map[string]any{
				"type":        "METRIC_KEY",
				"aggregation": "AVG",
				"metric_key":  "builtin:host.cpu.usage",
				"dimension_filter": []any{
					map[string]any{
						"filter": []any{
							map[string]any{
								"dimension_key":   "dt.entity.host",
								"dimension_value": "HOST-1234",
								"operator":        "EQUALS",
							},
						},
					},
				},
			},
		},
		{
			name: "MetricKey_WithPopulatedEntityFilter",
			input: map[string]any{
				"type":        "METRIC_KEY",
				"aggregation": "AVG",
				"metric_key":  "builtin:host.cpu.usage",
				"entity_filter": []any{
					map[string]any{
						"dimension_key": "dt.entity.host",
					},
				},
			},
		},
		{
			name: "MetricSelector_Minimal",
			input: map[string]any{
				"type":            "METRIC_SELECTOR",
				"metric_selector": "builtin:host.cpu.usage:avg",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			qd := unmarshalQueryDefinition(t, c.input)
			assert.NoError(t, qd.HandlePreconditions())
		})
	}
}

func TestQueryDefinition_HandlePreconditions_Invalid(t *testing.T) {
	cases := []struct {
		name        string
		input       map[string]any
		expectedErr string
	}{
		{
			name: "MetricSelector_WithAggregation",
			input: map[string]any{
				"type":            "METRIC_SELECTOR",
				"aggregation":     "AVG",
				"metric_selector": "builtin:host.cpu.usage:avg",
			},
			expectedErr: "'aggregation' must not be specified unless 'type' is set to 'METRIC_KEY'; got 'type'='METRIC_SELECTOR'",
		},
		{
			name: "MetricKey_WithoutAggregation",
			input: map[string]any{
				"type":       "METRIC_KEY",
				"metric_key": "builtin:host.cpu.usage",
			},
			expectedErr: "'aggregation' must be specified when 'type' is set to 'METRIC_KEY'; got 'type'='METRIC_KEY'",
		},
		{
			name: "MetricSelector_WithDimensionFilter",
			input: map[string]any{
				"type":            "METRIC_SELECTOR",
				"metric_selector": "builtin:host.cpu.usage:avg",
				"dimension_filter": []any{
					map[string]any{
						"filter": []any{
							map[string]any{
								"dimension_key":   "dt.entity.host",
								"dimension_value": "HOST-1234",
								"operator":        "EQUALS",
							},
						},
					},
				},
			},
			expectedErr: "'dimension_filter' must not be specified unless ('type' is set to 'METRIC_KEY' and 'metric_key' is set); got 'type'='METRIC_SELECTOR', 'metric_key'='<nil>'",
		},
		{
			name: "MetricKey_WithDimensionFilter_WithoutMetricKey",
			input: map[string]any{
				"type":        "METRIC_KEY",
				"aggregation": "AVG",
				"dimension_filter": []any{
					map[string]any{
						"filter": []any{
							map[string]any{
								"dimension_key":   "dt.entity.host",
								"dimension_value": "HOST-1234",
								"operator":        "EQUALS",
							},
						},
					},
				},
			},
			expectedErr: "'dimension_filter' must not be specified unless ('type' is set to 'METRIC_KEY' and 'metric_key' is set); got 'type'='METRIC_KEY', 'metric_key'='<nil>'",
		},
		{
			name: "MetricSelector_WithPopulatedEntityFilter",
			input: map[string]any{
				"type":            "METRIC_SELECTOR",
				"metric_selector": "builtin:host.cpu.usage:avg",
				"entity_filter": []any{
					map[string]any{
						"dimension_key": "dt.entity.host",
					},
				},
			},
			expectedErr: "'entity_filter' must not be specified unless ('type' is set to 'METRIC_KEY' and 'metric_key' is set); got 'type'='METRIC_SELECTOR', 'metric_key'='<nil>'",
		},
		{
			name: "MetricKey_WithEntityFilter_WithoutMetricKey",
			input: map[string]any{
				"type":        "METRIC_KEY",
				"aggregation": "AVG",
				"entity_filter": []any{
					map[string]any{
						"dimension_key": "dt.entity.host",
					},
				},
			},
			expectedErr: "'entity_filter' must not be specified unless ('type' is set to 'METRIC_KEY' and 'metric_key' is set); got 'type'='METRIC_KEY', 'metric_key'='<nil>'",
		},
		{
			name: "MetricSelector_WithMetricKey",
			input: map[string]any{
				"type":            "METRIC_SELECTOR",
				"metric_key":      "builtin:host.cpu.usage",
				"metric_selector": "builtin:host.cpu.usage:avg",
			},
			expectedErr: "'metric_key' must not be specified unless 'type' is set to 'METRIC_KEY'; got 'type'='METRIC_SELECTOR'",
		},
		{
			name: "MetricKey_WithoutMetricKey",
			input: map[string]any{
				"type":        "METRIC_KEY",
				"aggregation": "AVG",
			},
			expectedErr: "'metric_key' must be specified when 'type' is set to 'METRIC_KEY'; got 'type'='METRIC_KEY'",
		},
		{
			name: "MetricKey_WithMetricSelector",
			input: map[string]any{
				"type":            "METRIC_KEY",
				"aggregation":     "AVG",
				"metric_key":      "builtin:host.cpu.usage",
				"metric_selector": "builtin:host.cpu.usage:avg",
			},
			expectedErr: "'metric_selector' must not be specified unless 'type' is set to 'METRIC_SELECTOR'; got 'type'='METRIC_KEY'",
		},
		{
			name: "MetricSelector_WithoutMetricSelector",
			input: map[string]any{
				"type": "METRIC_SELECTOR",
			},
			expectedErr: "'metric_selector' must be specified when 'type' is set to 'METRIC_SELECTOR'; got 'type'='METRIC_SELECTOR'",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			qd := unmarshalQueryDefinition(t, c.input)
			assert.EqualError(t, qd.HandlePreconditions(), c.expectedErr)
		})
	}
}

// TestQueryDefinition_UnmarshalHCL_DropsStaleEmptyEntityFilter is the regression
// test for the bug where switching `type` from METRIC_KEY to METRIC_SELECTOR
// shipped an empty `entityFilter` to the API. An empty entity_filter can survive
// in state when (a) the resource was first created with type=METRIC_KEY without
// an entity_filter block (UnmarshalHCL auto-creates one), and (b) the
// DiffSuppressFunc on entity_filter masks the 1->0 removal during plan.
// UnmarshalHCL normalises that empty filter back to nil; this test pins the contract.
func TestQueryDefinition_UnmarshalHCL_DropsStaleEmptyEntityFilter(t *testing.T) {
	qd := unmarshalQueryDefinition(t, map[string]any{
		"type":            "METRIC_SELECTOR",
		"metric_selector": "builtin:host.cpu.usage:avg",
		// Empty entity_filter block — what survives in `d` when the
		// DiffSuppressFunc masks the removal of an auto-created filter.
		"entity_filter": []any{map[string]any{}},
	})
	assert.Nil(t, qd.EntityFilter, "empty entity_filter must be dropped when type=METRIC_SELECTOR; otherwise the API rejects the payload")
}

func unmarshalQueryDefinition(t *testing.T, input map[string]any) *QueryDefinition {
	t.Helper()
	s := new(QueryDefinition).Schema()
	d := schema.TestResourceDataRaw(t, s, input)
	assert.NotNil(t, d)

	qd := new(QueryDefinition)
	err := qd.UnmarshalHCL(hcl.DecoderFrom(d))
	assert.NoError(t, err)
	return qd
}
