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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/ingestsource/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var minimalIngestSource = settings.IngestSource{
	Kind:        "events",
	DisplayName: "displayName",
	PathSegment: "my.path.segment",
	Enabled:     true,
}

var bucket = "bucket"
var pipelineId = "pipelineId"
var sampleData = "sampleData"
var matcher = "not true"
var allFieldsSetIngestSource = settings.IngestSource{
	Kind:          "events",
	DefaultBucket: &bucket,
	DisplayName:   "displayName",
	Enabled:       false,
	PathSegment:   "some.path.segment",
	StaticRouting: &settings.PipelineReference{
		PipelineID:   &pipelineId,
		PipelineType: "custom",
	},
	Processing: &settings.Processing{
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
		input    settings.IngestSource
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
	s := new(settings.IngestSource).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected settings.IngestSource
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

			var actual settings.IngestSource
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
		input    settings.IngestSource
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

func TestIngestSource_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": "my.path.segment",
			},
			wantErr: false,
		},
		{
			name: "valid with all optional fields set",
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
			wantErr: false,
		},
		{
			name: "valid with digit",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": "my.0.1",
			},
			wantErr: false,
		},
		{
			name: "kind missing",
			input: map[string]interface{}{
				"display_name": "displayName",
				"path_segment": "my.path.segment",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "kind invalid",
			input: map[string]interface{}{
				"kind":         "invalid",
				"display_name": "displayName",
				"path_segment": "my.path.segment",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf("expected kind to be one of %q, got invalid", settings.AllowedKinds),
		},
		{
			name: "display_name missing",
			input: map[string]interface{}{
				"kind":         "events",
				"path_segment": "my.path.segment",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "display_name present but blank",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "",
				"path_segment": "my.path.segment",
			},
			wantErr: true,
			errMsg:  "expected length of display_name to be in the range (1 - 512), got ",
		},
		{
			name: "display_name too long",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": strings.Repeat("a", settings.DisplayNameMaxLength+1),
				"path_segment": "my.path.segment",
			},
			wantErr: true,
			errMsg:  "expected length of display_name to be in the range (1 - 512), got " + strings.Repeat("a", settings.DisplayNameMaxLength+1),
		},
		{
			name: "path_segment missing",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "path_segment present but blank",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": "",
			},
			wantErr: true,
			errMsg:  "expected length of path_segment to be in the range (1 - 100), got ",
		},
		{
			name: "path_segment present too long",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": strings.Repeat("a", settings.PathSegmentMaxLength+1),
			},
			wantErr: true,
			errMsg:  "expected length of path_segment to be in the range (1 - 100), got " + strings.Repeat("a", settings.PathSegmentMaxLength+1),
		},
		{
			name: "path_segment cannot start with number",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": "0abc",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf("invalid value for path_segment (%s)", settings.PathSegmentErrorMessage),
		},
		{
			name: "path_segment cannot start with dot",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": ".abc",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf("invalid value for path_segment (%s)", settings.PathSegmentErrorMessage),
		},
		{
			name: "path_segment cannot end with dot",
			input: map[string]interface{}{
				"kind":         "events",
				"display_name": "displayName",
				"path_segment": "abc.",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf("invalid value for path_segment (%s)", settings.PathSegmentErrorMessage),
		},
		{
			name: "default_bucket is present but blank",
			input: map[string]interface{}{
				"kind":           "events",
				"display_name":   "displayName",
				"path_segment":   "my.path.segment",
				"default_bucket": "",
			},
			wantErr: true,
			errMsg:  "expected length of default_bucket to be in the range (1 - 500), got ",
		},
		{
			name: "default_bucket too long",
			input: map[string]interface{}{
				"kind":           "events",
				"display_name":   "displayName",
				"path_segment":   "my.path.segment",
				"default_bucket": strings.Repeat("a", settings.DefaultBucketMaxLength+1),
			},
			wantErr: true,
			errMsg:  "expected length of default_bucket to be in the range (1 - 500), got " + strings.Repeat("a", settings.DefaultBucketMaxLength+1),
		},
	}

	r := &schema.Resource{Schema: new(settings.IngestSource).Schema()}
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
