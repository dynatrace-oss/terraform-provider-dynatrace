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

package processors_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleData = "my-sample-data"
var matcher = "not true"
var allFieldsSetDropProcessor = processors.Processor{
	Enabled:     true,
	Id:          "proc-1",
	Type:        processors.DropProcessorType,
	Description: "my-drop-processor",
	SampleData:  &sampleData,
	Matcher:     &matcher,
}

var minimalFieldsDropSetProcessor = processors.Processor{
	Enabled:     false,
	Id:          "proc-2",
	Type:        processors.DropProcessorType,
	Description: "my-other-drop-processor",
	Matcher:     &matcher,
}

var dqlProcessor = processors.Processor{
	Enabled:     true,
	Id:          "proc-2",
	Type:        processors.DqlProcessorType,
	Description: "my-proc-2",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	Dql:         &processors.DqlAttributes{Script: "fieldsAdd true"},
}

var fieldsAddProcessor = processors.Processor{
	Enabled:     true,
	Id:          "proc-3",
	Type:        processors.FieldsAddProcessorType,
	Description: "my-proc-3",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	FieldsAdd: &processors.FieldsAddAttributes{
		Fields: []*processors.FieldsAddItem{
			{
				Name:  "some-name",
				Value: "some-value",
			},
			{
				Name:  "some-other-name",
				Value: "some-other-value",
			},
		},
	},
}

var fieldsRenameProcessor = processors.Processor{
	Enabled:     true,
	Id:          "proc-4",
	Type:        processors.FieldsRenameProcessorType,
	Description: "my-proc-4",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	FieldsRename: &processors.FieldsRenameAttributes{
		Fields: []*processors.FieldsRenameItem{
			{
				FromName: "from-name",
				ToName:   "to-name",
			},
			{
				FromName: "from-other-name",
				ToName:   "to-other-name",
			},
		},
	},
}

var fieldsRemoveProcessor = processors.Processor{
	Enabled:     true,
	Id:          "proc-5",
	Type:        processors.FieldsRemoveProcessorType,
	Description: "my-proc-5",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	FieldsRemove: &processors.FieldsRemoveAttributes{
		Fields: []string{"to-remove-1", "to-remove-2"},
	},
}

func TestProcessor_MarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    processors.Processor
		expected []byte
	}{
		{
			name:  "drop-processor all fields set",
			input: allFieldsSetDropProcessor,
			expected: []byte(
				`{
				"id": "proc-1",
				"type": "drop",
				"matcher": "not true",
				"description": "my-drop-processor",
				"sampleData": "my-sample-data",
				"enabled": true
			}`),
		},
		{
			name:  "drop-processor minimal fields set",
			input: minimalFieldsDropSetProcessor,
			expected: []byte(
				`{
				"id": "proc-2",
				"type": "drop",
				"matcher": "not true",
				"description": "my-other-drop-processor",
				"enabled": false
			}`),
		},
		{
			name:  "dql-processor",
			input: dqlProcessor,
			expected: []byte(
				`{
				"id": "proc-2",
				"type": "dql",
				"matcher": "not true",
				"description": "my-proc-2",
				"sampleData": "my-sample-data",
				"enabled": true,
				"dql": {
					"script": "fieldsAdd true"
				}
			}`),
		},
		{
			name:  "fieldsAdd-processor",
			input: fieldsAddProcessor,
			expected: []byte(
				`{
				"id": "proc-3",
				"type": "fieldsAdd",
				"matcher": "not true",
				"description": "my-proc-3",
				"sampleData": "my-sample-data",
				"enabled": true,
				"fieldsAdd": {
					"fields": [
						{
							"name": "some-name",
							"value": "some-value"
						},
						{
							"name": "some-other-name",
							"value": "some-other-value"
						}
					]
				}
			}`),
		},
		{
			name:  "fieldsRename-processor",
			input: fieldsRenameProcessor,
			expected: []byte(
				`{
				"id": "proc-4",
				"type": "fieldsRename",
				"matcher": "not true",
				"description": "my-proc-4",
				"sampleData": "my-sample-data",
				"enabled": true,
				"fieldsRename": {
					"fields": [
						{
							"fromName": "from-name",
							"toName": "to-name"
						},
						{
							"fromName": "from-other-name",
							"toName": "to-other-name"
						}
					]
				}
			}`),
		},
		{
			name:  "fieldsRemove-processor",
			input: fieldsRemoveProcessor,
			expected: []byte(
				`{
				"id": "proc-5",
				"type": "fieldsRemove",
				"matcher": "not true",
				"description": "my-proc-5",
				"sampleData": "my-sample-data",
				"enabled": true,
				"fieldsRemove": {
					"fields": [
						"to-remove-1",
						"to-remove-2"
					]
				}
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

func TestProcessor_MarshalHCL(t *testing.T) {
	cases := []struct {
		name     string
		input    processors.Processor
		expected hcl.Properties
	}{
		{
			name:  "drop-processor all fields set",
			input: allFieldsSetDropProcessor,
			expected: hcl.Properties{
				"enabled":     true,
				"type":        "drop",
				"id":          "proc-1",
				"sample_data": sampleData,
				"description": "my-drop-processor",
				"matcher":     "not true",
			},
		},
		{
			name:  "drop-processor minimal fields set",
			input: minimalFieldsDropSetProcessor,
			expected: hcl.Properties{
				"enabled":     false,
				"type":        "drop",
				"id":          "proc-2",
				"description": "my-other-drop-processor",
				"matcher":     "not true",
			},
		},
		{
			name:  "dql-processor",
			input: dqlProcessor,
			expected: hcl.Properties{
				"enabled":     true,
				"type":        "dql",
				"id":          "proc-2",
				"description": "my-proc-2",
				"sample_data": sampleData,
				"matcher":     "not true",
				"dql": []interface{}{
					hcl.Properties{
						"script": "fieldsAdd true",
					},
				},
			},
		},
		{
			name:  "fieldsAdd-processor",
			input: fieldsAddProcessor,
			expected: hcl.Properties{
				"enabled":     true,
				"type":        "fieldsAdd",
				"id":          "proc-3",
				"description": "my-proc-3",
				"sample_data": sampleData,
				"matcher":     "not true",
				"fields_add": []interface{}{
					hcl.Properties{
						"field": []interface{}{
							hcl.Properties{
								"name":  "some-name",
								"value": "some-value",
							},
							hcl.Properties{
								"name":  "some-other-name",
								"value": "some-other-value",
							},
						},
					},
				},
			},
		},
		{
			name:  "fieldsRename-processor",
			input: fieldsRenameProcessor,
			expected: hcl.Properties{
				"enabled":     true,
				"type":        "fieldsRename",
				"id":          "proc-4",
				"description": "my-proc-4",
				"sample_data": sampleData,
				"matcher":     "not true",
				"fields_rename": []interface{}{
					hcl.Properties{
						"field": []interface{}{
							hcl.Properties{
								"from_name": "from-name",
								"to_name":   "to-name",
							},
							hcl.Properties{
								"from_name": "from-other-name",
								"to_name":   "to-other-name",
							},
						},
					},
				},
			},
		},
		{
			name:  "fieldsRemove-processor",
			input: fieldsRemoveProcessor,
			expected: hcl.Properties{
				"enabled":     true,
				"type":        "fieldsRemove",
				"id":          "proc-5",
				"description": "my-proc-5",
				"sample_data": sampleData,
				"matcher":     "not true",
				"fields_remove": []interface{}{
					hcl.Properties{
						"fields": []string{
							"to-remove-1",
							"to-remove-2",
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

func TestProcessor_UnmarshalHCL(t *testing.T) {
	s := new(processors.Processor).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected processors.Processor
	}{
		{
			name: "all fields set",
			input: map[string]interface{}{
				"enabled":     true,
				"type":        "drop",
				"id":          "proc-1",
				"sample_data": sampleData,
				"description": "my-drop-processor",
				"matcher":     "not true",
			},
			expected: allFieldsSetDropProcessor,
		},
		{
			name: "minimal fields set",
			input: map[string]interface{}{
				"enabled":     false,
				"type":        "drop",
				"id":          "proc-2",
				"description": "my-other-drop-processor",
				"matcher":     "not true",
			},
			expected: minimalFieldsDropSetProcessor,
		},
		{
			name: "dql-processor",
			input: map[string]interface{}{
				"enabled":     true,
				"type":        "dql",
				"id":          "proc-2",
				"description": "my-proc-2",
				"matcher":     "not true",
				"sample_data": sampleData,
				"dql": []interface{}{
					map[string]interface{}{
						"script": "fieldsAdd true",
					},
				},
			},
			expected: dqlProcessor,
		},
		{
			name: "fieldsAdd-processor",
			input: map[string]interface{}{
				"enabled":     true,
				"type":        "fieldsAdd",
				"id":          "proc-3",
				"description": "my-proc-3",
				"matcher":     "not true",
				"sample_data": sampleData,
				"fields_add": []interface{}{
					map[string]interface{}{
						"field": []interface{}{
							map[string]interface{}{
								"name":  "some-name",
								"value": "some-value",
							},
							map[string]interface{}{
								"name":  "some-other-name",
								"value": "some-other-value",
							},
						},
					},
				},
			},
			expected: fieldsAddProcessor,
		},
		{
			name: "fieldsRename-processor",
			input: map[string]interface{}{
				"enabled":     true,
				"type":        "fieldsRename",
				"id":          "proc-4",
				"description": "my-proc-4",
				"matcher":     "not true",
				"sample_data": sampleData,
				"fields_rename": []interface{}{
					map[string]interface{}{
						"field": []interface{}{
							map[string]interface{}{
								"from_name": "from-name",
								"to_name":   "to-name",
							},
							map[string]interface{}{
								"from_name": "from-other-name",
								"to_name":   "to-other-name",
							},
						},
					},
				},
			},
			expected: fieldsRenameProcessor,
		},
		{
			name: "fieldsRemove-processor",
			input: map[string]interface{}{
				"enabled":     true,
				"type":        "fieldsRemove",
				"id":          "proc-5",
				"description": "my-proc-5",
				"matcher":     "not true",
				"sample_data": sampleData,
				"fields_remove": []interface{}{
					map[string]interface{}{
						"fields": []interface{}{
							"to-remove-1",
							"to-remove-2",
						},
					},
				},
			},
			expected: fieldsRemoveProcessor,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual processors.Processor
			decoder := hcl.DecoderFrom(d)
			err := actual.UnmarshalHCL(decoder)
			assert.Equal(t, c.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestProcessor_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "proc-1",
				"description": "my-drop-processor",
			},
			wantErr: false,
		},
		{
			name: "type missing",
			input: map[string]interface{}{
				"id":          "proc-1",
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "invalid processor type",
			input: map[string]interface{}{
				"type":        "invalid",
				"id":          "proc-1",
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  fmt.Sprintf("expected type to be one of %q, got invalid", processors.AvailableProcessorTypes),
		},
		{
			name: "id missing",
			input: map[string]interface{}{
				"type":        "drop",
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "id too short",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          strings.Repeat("a", processors.IDMinLength-1),
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  "expected length of id to be in the range (4 - 100), got aaa",
		},
		{
			name: "id too long",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          strings.Repeat("a", processors.IDMaxLength+1),
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  "expected length of id to be in the range (4 - 100), got " + strings.Repeat("a", processors.IDMaxLength+1),
		},
		{
			name: "id starts with dt.",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "dt.2",
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  "id must not start with 'dt.' or 'dynatrace.'",
		},
		{
			name: "id starts with dynatrace.",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "dynatrace.2",
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  "id must not start with 'dt.' or 'dynatrace.'",
		},
		{
			name: "description missing",
			input: map[string]interface{}{
				"type": "drop",
				"id":   "proc-1",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "description present but blank",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "processor1",
				"description": "",
			},
			wantErr: true,
			errMsg:  "expected length of description to be in the range (1 - 512), got ",
		},
		{
			name: "description too long",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "processor1",
				"description": strings.Repeat("a", processors.DescriptionMaxLength+1),
			},
			wantErr: true,
			errMsg:  "expected length of description to be in the range (1 - 512), got " + strings.Repeat("a", processors.DescriptionMaxLength+1),
		},
		{
			name: "matcher present but blank",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "processor1",
				"description": "my-drop-processor",
				"matcher":     "",
			},
			wantErr: true,
			errMsg:  "expected length of matcher to be in the range (1 - 1500), got ",
		},
		{
			name: "matcher too long",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "processor1",
				"description": "my-drop-processor",
				"matcher":     strings.Repeat("a", processors.MatcherMaxLength+1),
			},
			wantErr: true,
			errMsg:  "expected length of matcher to be in the range (1 - 1500), got " + strings.Repeat("a", processors.MatcherMaxLength+1),
		},
		{
			name: "sample_data present but blank",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "processor1",
				"description": "my-drop-processor",
				"sample_data": "",
			},
			wantErr: true,
			errMsg:  "expected length of sample_data to be in the range (1 - 8192), got ",
		},
		{
			name: "sample_data too long",
			input: map[string]interface{}{
				"type":        "drop",
				"id":          "processor1",
				"description": "my-drop-processor",
				"sample_data": strings.Repeat("a", processors.SampleDataMaxLength+1),
			},
			wantErr: true,
			errMsg:  "expected length of sample_data to be in the range (1 - 8192), got " + strings.Repeat("a", processors.SampleDataMaxLength+1),
		},
	}

	r := &schema.Resource{Schema: new(processors.Processor).Schema()}
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
