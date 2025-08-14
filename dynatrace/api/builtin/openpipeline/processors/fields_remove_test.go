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
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestFieldsRemoveAttributes_MarshalHCL(t *testing.T) {
	var input = FieldsRemoveAttributes{
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
	var expected = FieldsRemoveAttributes{Fields: []string{"field1", "field2"}}

	d := schema.TestResourceDataRaw(t, new(FieldsRemoveAttributes).Schema(), input)
	assert.NotNil(t, d)

	var actual FieldsRemoveAttributes
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
			name:    "empty",
			input:   map[string]interface{}{},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name:    "empty-fields",
			input:   map[string]interface{}{"fields": []interface{}{}},
			wantErr: false,
		},
		{
			name:    "unrecognized-property",
			input:   map[string]interface{}{"fields": []interface{}{}, "strange-property": true},
			wantErr: true,
			errMsg:  "Invalid or unknown key",
		},
	}

	r := &schema.Resource{Schema: new(FieldsRemoveAttributes).Schema()}
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
