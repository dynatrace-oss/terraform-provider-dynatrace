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

func TestFieldsRenameAttributes_MarshalHCL(t *testing.T) {
	var input = processors.FieldsRenameAttributes{
		Fields: []*processors.FieldsRenameItem{
			{
				FromName: "some-name",
				ToName:   "new-name",
			},
			{
				FromName: "some-other-name",
				ToName:   "other-new-name",
			},
		},
	}
	var expected = hcl.Properties{
		"field": []interface{}{
			hcl.Properties{
				"from_name": "some-name",
				"to_name":   "new-name",
			},
			hcl.Properties{
				"from_name": "some-other-name",
				"to_name":   "other-new-name",
			},
		},
	}

	var actual = hcl.Properties{}

	err := input.MarshalHCL(actual)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestFieldsRenameAttributes_UnmarshalHCL(t *testing.T) {
	var input = map[string]interface{}{
		"field": []interface{}{
			map[string]interface{}{
				"from_name": "some-name",
				"to_name":   "new-name",
			},
			map[string]interface{}{
				"from_name": "some-other-name",
				"to_name":   "other-new-name",
			},
		},
	}
	var expected = processors.FieldsRenameAttributes{
		Fields: []*processors.FieldsRenameItem{
			{
				FromName: "some-name",
				ToName:   "new-name",
			},
			{
				FromName: "some-other-name",
				ToName:   "other-new-name",
			},
		},
	}

	d := schema.TestResourceDataRaw(t, new(processors.FieldsRenameAttributes).Schema(), input)
	assert.NotNil(t, d)

	var actual processors.FieldsRenameAttributes
	decoder := hcl.DecoderFrom(d)

	err := actual.UnmarshalHCL(decoder)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestFieldsRenameAttributes_Validate(t *testing.T) {
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
						"from_name": "some-name",
						"to_name":   "new-name",
					},
					map[string]interface{}{
						"from_name": "some-other-name",
						"to_name":   "other-new-name",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty from_name and to_name is valid",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"from_name": "",
						"to_name":   "",
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
			name: "field-missing-from-name",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"to_name": "new-name",
					},
				},
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "field-missing-to-name",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"from_name": "some-name",
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
						"from_name":        "some-name",
						"to_name":          "new-name",
						"strange-property": true,
					},
				},
			},
			wantErr: true,
			errMsg:  "Invalid or unknown key",
		},
		{
			name: "from-name too long",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"from_name": strings.Repeat("a", processors.FieldsRenameMaxNameLength+1),
						"to_name":   "new-name",
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of field.0.from_name to be in the range (0 - 256), got " + strings.Repeat("a", processors.FieldsRenameMaxNameLength+1),
		},
		{
			name: "to-name too long",
			input: map[string]interface{}{
				"field": []interface{}{
					map[string]interface{}{
						"from_name": "old-name",
						"to_name":   strings.Repeat("a", processors.FieldsRenameMaxNameLength+1),
					},
				},
			},
			wantErr: true,
			errMsg:  "expected length of field.0.to_name to be in the range (0 - 256), got " + strings.Repeat("a", processors.FieldsRenameMaxNameLength+1),
		},
		{
			name: "too many fields",
			input: map[string]interface{}{
				"field": createNFieldsRenameFields(51),
			},
			wantErr: true,
			errMsg:  "Too many list items",
		},
	}

	r := &schema.Resource{Schema: new(processors.FieldsRenameAttributes).Schema()}
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

func createNFieldsRenameFields(n int) []interface{} {
	var fields []interface{}
	var field = map[string]interface{}{
		"from_name": "old-name",
		"to_name":   "new-name",
	}

	for i := 0; i < n; i++ {
		fields = append(fields, field)
	}

	return fields
}
