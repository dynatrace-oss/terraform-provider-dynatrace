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

func TestDqlAttributes_MarshalHCL(t *testing.T) {
	var input = DqlAttributes{Script: "fieldsAdd true"}
	var actual = hcl.Properties{}

	err := input.MarshalHCL(actual)

	assert.Equal(t, hcl.Properties{"script": "fieldsAdd true"}, actual)
	assert.NoError(t, err)
}

func TestDqlAttributes_UnmarshalHCL(t *testing.T) {
	d := schema.TestResourceDataRaw(t,
		new(DqlAttributes).Schema(),
		map[string]interface{}{"script": "fieldsAdd true"})
	assert.NotNil(t, d)
	var actual DqlAttributes
	decoder := hcl.DecoderFrom(d)

	err := actual.UnmarshalHCL(decoder)

	assert.Equal(t, DqlAttributes{Script: "fieldsAdd true"}, actual)
	assert.NoError(t, err)
}

func TestDqlAttributes_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			input: map[string]interface{}{
				"script": "fieldsAdd true",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			input: map[string]interface{}{
				"script":        "fieldsAdd true",
				"unknown-field": true,
			},
			wantErr: true,
			errMsg:  "Invalid or unknown key",
		},
		{
			name:    "empty",
			input:   map[string]interface{}{},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
	}

	r := &schema.Resource{Schema: new(DqlAttributes).Schema()}
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
