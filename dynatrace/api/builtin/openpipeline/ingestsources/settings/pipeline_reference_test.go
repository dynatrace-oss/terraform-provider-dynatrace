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
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestPipelineReference_MarshalHCL(t *testing.T) {
	var pId = "pipelineId"
	cases := []struct {
		name     string
		input    PipelineReference
		expected hcl.Properties
	}{
		{
			name: "builtin-pipeline-reference",
			input: PipelineReference{
				BuiltinPipelineID: &pId,
				PipelineType:      "builtin",
			},
			expected: hcl.Properties{
				"pipeline_id":   pId,
				"pipeline_type": "custom",
			},
		},
		{
			name: "custom-pipeline-reference",
			input: PipelineReference{
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
	s := new(PipelineReference).Schema()
	var pId = "pipelineId"

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected PipelineReference
	}{
		{
			name: "builtin pipeline reference",
			input: map[string]interface{}{
				"builtin_pipeline_id": pId,
				"pipeline_type":       "custom",
			},
			expected: PipelineReference{
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
			expected: PipelineReference{
				PipelineID:   &pId,
				PipelineType: "custom",
			},
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
