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
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestFieldsAddAttributes_MarshalHCL(t *testing.T) {
	var input = processors.FieldsAddAttributes{
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
	}
	var expected = hcl.Properties{
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
	}

	var actual = hcl.Properties{}

	err := input.MarshalHCL(actual)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestFieldsAddAttributes_UnmarshalHCL(t *testing.T) {
	var input = map[string]interface{}{
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
	}
	var expected = processors.FieldsAddAttributes{
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
	}

	d := schema.TestResourceDataRaw(t, new(processors.FieldsAddAttributes).Schema(), input)
	assert.NotNil(t, d)

	var actual processors.FieldsAddAttributes
	decoder := hcl.DecoderFrom(d)

	err := actual.UnmarshalHCL(decoder)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestFieldsAddAttributes_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			input: map[string]interface{}{
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
			wantErr: false,
		},
		{
			name:    "empty",
			input:   map[string]interface{}{},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "empty-fields",
			input: map[string]interface{}{
				"field": []interface{}{},
			},
			wantErr: true,
			errMsg:  "Not enough list items",
		},
		{
			name: "empty-field",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{},
				},
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "field-missing-name",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"value": "some-value",
					},
				},
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "field-missing-value",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"name": "some-name",
					},
				},
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "field-with-unrecognized-property",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"name":             "some-name",
						"value":            "some-value",
						"strange-property": true,
					},
				},
			},
			wantErr: true,
			errMsg:  "Invalid or unknown key",
		},
		{
			name: "name and value empty is valid",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"name":  "",
						"value": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "field name too long",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"name":  strings.Repeat("a", processors.FieldsAddMaxNameLength+1),
						"value": "some-value",
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of field.0.name to be in the range (0 - 256), got " + strings.Repeat("a", processors.FieldsAddMaxNameLength+1),
		},
		{
			name: "field value too long",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"name":  "some-name",
						"value": strings.Repeat("a", processors.FieldsAddMaxValueLength+1),
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of field.0.value to be in the range (0 - 512), got " + strings.Repeat("a", processors.FieldsAddMaxValueLength+1),
		},
		{
			name: "too-many-fields",
			input: map[string]interface{}{
				"field": createNFieldAddFields(51),
			},
			wantErr: true,
			errMsg:  "Too many list items",
		},
	}

	r := &schema.Resource{Schema: new(processors.FieldsAddAttributes).Schema()}
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

func createNFieldAddFields(n int) []interface{} {
	var fields []interface{}
	var field = map[string]interface{}{
		"name":  "some-name",
		"value": "some-value",
	}

	for i := 0; i < n; i++ {
		fields = append(fields, field)
	}

	return fields
}
