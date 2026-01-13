//go:build unit

/**
* @license
* Copyright 2026 Dynatrace LLC
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

package comparisoninfo

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseComparisonInfo_UnmarshalHCL_MarshalJSON(t *testing.T) {
	s := new(BaseComparisonInfo).Schema()

	cases := []struct {
		name         string
		input        map[string]interface{}
		expectedJSON string
	}{
		{
			name: "basic type only - STRING",
			input: map[string]interface{}{
				"type": "STRING",
			},
			expectedJSON: `{"negate":false,"type":"STRING"}`,
		},
		{
			name: "empty unknowns string",
			input: map[string]interface{}{
				"type":     "STRING",
				"unknowns": `{}`,
			},
			expectedJSON: `{"negate":false,"type":"STRING"}`,
		},
		{
			name: "basic type only - NUMBER",
			input: map[string]interface{}{
				"type": "NUMBER",
			},
			expectedJSON: `{"negate":false,"type":"NUMBER"}`,
		},
		{
			name: "basic type only - BOOLEAN",
			input: map[string]interface{}{
				"type": "BOOLEAN",
			},
			expectedJSON: `{"negate":false,"type":"BOOLEAN"}`,
		},
		{
			name: "type with unknowns containing extra fields",
			input: map[string]interface{}{
				"type":     "STRING",
				"unknowns": `{"comparison":"EQUALS","value":"test-value"}`,
			},
			expectedJSON: `{"negate":false,"type":"STRING","comparison":"EQUALS","value":"test-value"}`,
		},
		{
			name: "unknowns with negate set to true",
			input: map[string]interface{}{
				"type":     "STRING",
				"unknowns": `{"negate":true,"comparison":"CONTAINS"}`,
			},
			// negate from unknowns is unmarshalled, then MarshalJSON outputs it
			expectedJSON: `{"negate":true,"type":"STRING","comparison":"CONTAINS"}`,
		},
		{
			name: "unknowns containing type field - explicit type takes precedence",
			input: map[string]interface{}{
				"type":     "NUMBER",
				"unknowns": `{"type":"STRING","comparison":"EQUALS"}`,
			},
			// The explicit "type" field in HCL takes precedence over unknowns
			expectedJSON: `{"negate":false,"type":"NUMBER","comparison":"EQUALS"}`,
		},
		{
			name: "unknowns with nested object",
			input: map[string]interface{}{
				"type":     "TAG",
				"unknowns": `{"value":{"context":"CONTEXTLESS","key":"my-tag"}}`,
			},
			expectedJSON: `{"negate":false,"type":"TAG","value":{"context":"CONTEXTLESS","key":"my-tag"}}`,
		},
		{
			name: "unknowns allows for complex type - STRING_ONE_AGENT_ATTRIBUTE",
			input: map[string]interface{}{
				"type":     "STRING_ONE_AGENT_ATTRIBUTE",
				"unknowns": `{"caseSensitive":false,"comparison":"CONTAINS","oneAgentAttributeKey":"http.route","value":"/services/*"}`,
			},
			expectedJSON: `{"negate":false,"type":"STRING_ONE_AGENT_ATTRIBUTE","caseSensitive":false,"comparison":"CONTAINS","oneAgentAttributeKey":"http.route","value":"/services/*"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Create test resource data from HCL input
			d := schema.TestResourceDataRaw(t, s, c.input)
			require.NotNil(t, d)

			// Unmarshal HCL into BaseComparisonInfo
			var info BaseComparisonInfo
			decoder := hcl.DecoderFrom(d)
			err := info.UnmarshalHCL(decoder)
			require.NoError(t, err)

			// Marshal to JSON
			actualJSON, err := info.MarshalJSON()
			require.NoError(t, err)

			// Compare JSON output
			assert.JSONEq(t, c.expectedJSON, string(actualJSON))
		})
	}
}
