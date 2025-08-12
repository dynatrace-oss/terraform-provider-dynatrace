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
	"strconv"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestFieldsRemoveAttributes_MarshalHCL(t *testing.T) {
	var input = processors.FieldsRemoveAttributes{
		Fields: []string{"field1", "field2"},
	}
	var expected = hcl.Properties{
		"fields": []string{"field1", "field2"},
	}

	var actual = hcl.Properties{}

	err := input.MarshalHCL(actual)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestFieldsRemoveAttributes_UnmarshalHCL(t *testing.T) {
	var input = map[string]interface{}{
		"fields": []interface{}{"field1", "field2"},
	}
	var expected = processors.FieldsRemoveAttributes{Fields: []string{"field1", "field2"}}

	d := schema.TestResourceDataRaw(t, new(processors.FieldsRemoveAttributes).Schema(), input)
	assert.NotNil(t, d)

	var actual processors.FieldsRemoveAttributes
	decoder := hcl.DecoderFrom(d)

	err := actual.UnmarshalHCL(decoder)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestFieldsRemoveAttributes_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid",
			input:   map[string]interface{}{"fields": []interface{}{"field1", "field2"}},
			wantErr: false,
		},
		{
			name:    "duplicate fields, still valid",
			input:   map[string]interface{}{"fields": []interface{}{"field1", "field1"}},
			wantErr: false,
		},
		{
			name:    "empty fields are also valid",
			input:   map[string]interface{}{"fields": []interface{}{"", ""}},
			wantErr: false,
		},
		{
			name:    "empty",
			input:   map[string]interface{}{},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name:    "empty-fields",
			input:   map[string]interface{}{"fields": []interface{}{}},
			wantErr: true,
			errMsg:  "Not enough list items",
		},
		{
			name:    "unrecognized-property",
			input:   map[string]interface{}{"fields": []interface{}{"field1"}, "strange-property": true},
			wantErr: true,
			errMsg:  "Invalid or unknown key",
		},
		{
			name:    "too many fields",
			input:   map[string]interface{}{"fields": createNFieldsRemoveFields(51)},
			wantErr: true,
			errMsg:  "Too many list items",
		},
	}

	r := &schema.Resource{Schema: new(processors.FieldsRemoveAttributes).Schema()}
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

func createNFieldsRemoveFields(n int) []interface{} {
	var fields []interface{}
	for i := 0; i < n; i++ {
		fields = append(fields, "field"+strconv.Itoa(i))
	}
	return fields
}
