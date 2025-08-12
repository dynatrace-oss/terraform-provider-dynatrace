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

package processors

import (
	"encoding/json"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleData = "my-sample-data"
var matcher = "not true"
var allFieldsSetDropProcessor = Processor{
	Enabled:     true,
	Id:          "proc-1",
	Type:        DropProcessorType,
	Description: "my-drop-processor",
	SampleData:  &sampleData,
	Matcher:     &matcher,
}

var minimalFieldsDropSetProcessor = Processor{
	Enabled:     false,
	Id:          "proc-2",
	Type:        DropProcessorType,
	Description: "my-other-drop-processor",
	Matcher:     &matcher,
}

var dqlProcessor = Processor{
	Enabled:     true,
	Id:          "proc-2",
	Type:        DqlProcessorType,
	Description: "my-proc-2",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	Dql:         &DqlAttributes{Script: "fieldsAdd true"},
}

var fieldsAddProcessor = Processor{
	Enabled:     true,
	Id:          "proc-3",
	Type:        FieldsAddProcessorType,
	Description: "my-proc-3",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	FieldsAdd: &FieldsAddAttributes{
		Fields: []*FieldsAddItem{
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

var fieldsRenameProcessor = Processor{
	Enabled:     true,
	Id:          "proc-4",
	Type:        FieldsRenameProcessorType,
	Description: "my-proc-4",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	FieldsRename: &FieldsRenameAttributes{
		Fields: []*FieldsRenameItem{
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

var fieldsRemoveProcessor = Processor{
	Enabled:     true,
	Id:          "proc-5",
	Type:        FieldsRemoveProcessorType,
	Description: "my-proc-5",
	SampleData:  &sampleData,
	Matcher:     &matcher,
	FieldsRemove: &FieldsRemoveAttributes{
		Fields: []string{"to-remove-1", "to-remove-2"},
	},
}

func TestProcessor_MarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    Processor
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
		input    Processor
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
				"fieldsAdd": []interface{}{
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
				"fieldsRename": []interface{}{
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
				"fieldsRemove": []interface{}{
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
	s := new(Processor).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected Processor
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
				"fieldsAdd": []interface{}{
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
				"fieldsRename": []interface{}{
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
				"fieldsRemove": []interface{}{
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

			var actual Processor
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
			errMsg:  "The argument \"type\" is required, but no definition was found.",
		},
		{
			name: "id missing",
			input: map[string]interface{}{
				"type":        "drop",
				"description": "my-drop-processor",
			},
			wantErr: true,
			errMsg:  "The argument \"id\" is required, but no definition was found.",
		},
		{
			name: "description missing",
			input: map[string]interface{}{
				"type": "drop",
				"id":   "proc-1",
			},
			wantErr: true,
			errMsg:  "The argument \"description\" is required, but no definition was found.",
		},
	}

	r := &schema.Resource{Schema: new(Processor).Schema()}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := terraform.NewResourceConfigRaw(tt.input)

			diags := r.Validate(c)
			assert.Equal(t, tt.wantErr, diags.HasError())

			if diags.HasError() {
				assert.Equal(t, tt.errMsg, diags[0].Detail)
			}
		})
	}
}
