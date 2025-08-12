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
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/ingestsource/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestPipelineReference_MarshalHCL(t *testing.T) {
	var pId = "pipelineId"
	cases := []struct {
		name     string
		input    settings.PipelineReference
		expected hcl.Properties
	}{
		{
			name: "builtin-pipeline-reference",
			input: settings.PipelineReference{
				BuiltinPipelineID: &pId,
				PipelineType:      "builtin",
			},
			expected: hcl.Properties{
				"builtin_pipeline_id": pId,
				"pipeline_type":       "builtin",
			},
		},
		{
			name: "custom-pipeline-reference",
			input: settings.PipelineReference{
				PipelineID:   &pId,
				PipelineType: "custom",
			},
			expected: hcl.Properties{
				"pipeline_id":   pId,
				"pipeline_type": "custom",
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

func TestPipelineReference_UnmarshalHCL(t *testing.T) {
	s := new(settings.PipelineReference).Schema("")
	var pId = "pipelineId"

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected settings.PipelineReference
	}{
		{
			name: "builtin pipeline reference",
			input: map[string]interface{}{
				"builtin_pipeline_id": pId,
				"pipeline_type":       "builtin",
			},
			expected: settings.PipelineReference{
				BuiltinPipelineID: &pId,
				PipelineType:      "builtin",
			},
		},
		{
			name: "custom pipeline reference",
			input: map[string]interface{}{
				"pipeline_id":   pId,
				"pipeline_type": "custom",
			},
			expected: settings.PipelineReference{
				PipelineID:   &pId,
				PipelineType: "custom",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual settings.PipelineReference
			decoder := hcl.DecoderFrom(d)
			err := actual.UnmarshalHCL(decoder)
			assert.Equal(t, c.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestPipelineReference_Validate(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			input: map[string]interface{}{
				"builtin_pipeline_id": "pId",
				"pipeline_type":       "builtin",
			},
			wantErr: false,
		},
		{
			name: "type missing",
			input: map[string]interface{}{
				"builtin_pipeline_id": "pId",
			},
			wantErr: true,
			errMsg:  "Missing required argument",
		},
		{
			name: "no id set",
			input: map[string]interface{}{
				"pipeline_type": "builtin",
			},
			wantErr: true,
			errMsg:  "Invalid combination of arguments",
		},
		{
			name: "both ids set",
			input: map[string]interface{}{
				"builtin_pipeline_id": "pId",
				"pipeline_id":         "pId",
				"pipeline_type":       "builtin",
			},
			wantErr: true,
			errMsg:  "Invalid combination of arguments",
		},
		{
			name: "type invalid",
			input: map[string]interface{}{
				"builtin_pipeline_id": "pId",
				"pipeline_type":       "invalid",
			},
			wantErr: true,
			errMsg:  "expected pipeline_type to be one of [\"custom\" \"builtin\"], got invalid",
		},
		{
			name: "builtin_pipeline_id present but blank",
			input: map[string]interface{}{
				"builtin_pipeline_id": "",
				"pipeline_type":       "builtin",
			},
			wantErr: true,
			errMsg:  "expected length of builtin_pipeline_id to be in the range (1 - 500), got ",
		},
		{
			name: "builtin_pipeline_id too long",
			input: map[string]interface{}{
				"builtin_pipeline_id": strings.Repeat("a", settings.BuiltinPipelineIDMaxLength+1),
				"pipeline_type":       "builtin",
			},
			wantErr: true,
			errMsg:  "expected length of builtin_pipeline_id to be in the range (1 - 500), got " + strings.Repeat("a", settings.BuiltinPipelineIDMaxLength+1),
		},
	}

	r := &schema.Resource{Schema: new(settings.PipelineReference).Schema("")}
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
