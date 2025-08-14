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

package settings

import (
	"encoding/json"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var minimalIngestSource = IngestSource{
	Kind:        "events",
	DisplayName: "displayName",
	Enabled:     true,
	PathSegment: "my.path.segment",
}

var bucket = "bucket"
var pipelineId = "pipelineId"
var sampleData = "sampleData"
var matcher = "not true"
var allFieldsSetIngestSource = IngestSource{
	Kind:          "events",
	DefaultBucket: &bucket,
	DisplayName:   "displayName",
	Enabled:       false,
	PathSegment:   "some.path.segment",
	StaticRouting: &PipelineReference{
		PipelineID:   &pipelineId,
		PipelineType: "custom",
	},
	Processing: &Processing{
		Processors: []*processors.Processor{
			{
				Enabled:     true,
				Id:          "proc-2",
				Type:        processors.DqlProcessorType,
				Description: "my-proc-2",
				SampleData:  &sampleData,
				Matcher:     &matcher,
				Dql:         &processors.DqlAttributes{Script: "fieldsAdd true"},
			},
		},
	},
}

func TestIngestSource_MarshalHCL(t *testing.T) {
	cases := []struct {
		name     string
		input    IngestSource
		expected hcl.Properties
	}{
		{
			name:  "minimum-set",
			input: minimalIngestSource,
			expected: hcl.Properties{
				"kind":         "events",
				"display_name": "displayName",
				"enabled":      true,
				"path_segment": "my.path.segment",
			},
		},
		{
			name:  "all-set",
			input: allFieldsSetIngestSource,
			expected: hcl.Properties{
				"kind":           "events",
				"default_bucket": "bucket",
				"display_name":   "displayName",
				"enabled":        false,
				"path_segment":   "some.path.segment",
				"static_routing": []interface{}{
					hcl.Properties{
						"pipeline_id":   "pipelineId",
						"pipeline_type": "custom",
					},
				},
				"processing": []interface{}{
					hcl.Properties{
						"processor": []interface{}{
							hcl.Properties{
								"enabled":     true,
								"id":          "proc-2",
								"type":        processors.DqlProcessorType,
								"description": "my-proc-2",
								"sample_data": sampleData,
								"matcher":     matcher,
								"dql": []interface{}{
									hcl.Properties{
										"script": "fieldsAdd true",
									},
								},
							},
						},
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

func TestIngestSource_UnmarshalHCL(t *testing.T) {
	s := new(IngestSource).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected IngestSource
	}{
		{
			name: "minimal fields set",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": "my.path.segment",
			},
			expected: minimalIngestSource,
		},
		{
			name: "all fields set",
			input: map[string]interface{}{
				"kind":           "events",
				"default_bucket": "bucket",
				"display_name":   "displayName",
				"path_segment":   "some.path.segment",
				"enabled":        false,
				"static_routing": []interface{}{
					map[string]interface{}{
						"pipeline_id":   "pipelineId",
						"pipeline_type": "custom",
					},
				},
				"processing": []interface{}{
					map[string]interface{}{
						"processor": []interface{}{
							map[string]interface{}{
								"id":          "proc-2",
								"type":        "dql",
								"matcher":     "not true",
								"description": "my-proc-2",
								"sample_data": "sampleData",
								"enabled":     true,
								"dql": []interface{}{
									map[string]interface{}{
										"script": "fieldsAdd true",
									},
								},
							},
						},
					},
				},
			},
			expected: allFieldsSetIngestSource,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual IngestSource
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
		input    IngestSource
		expected []byte
	}{
		{
			name:  "minimal-fields-set",
			input: minimalIngestSource,
			expected: []byte( // The API demands that "processing" is present, even if it is null, so that seems correct.
				`{
					"displayName": "displayName",
					"enabled": true,
					"pathSegment": "my.path.segment",
					"processing": {}
				}`),
		},
		{
			name:  "all-fields-set",
			input: allFieldsSetIngestSource,
			expected: []byte(
				`{
					"defaultBucket": "bucket",	
					"displayName": "displayName",
					"pathSegment": "some.path.segment",
					"enabled": false,
					"staticRouting": {
					  "pipelineId": "pipelineId",
					  "pipelineType": "custom"
					},
					"processing": {
						"processors": [
							{
								"id": "proc-2",
								"type": "dql",
								"matcher": "not true",
								"description": "my-proc-2",
								"sampleData": "sampleData",
								"enabled": true,
								"dql": {
									"script": "fieldsAdd true"
								}
							}
						]
					}	
				}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.input.MarshalJSON()
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
