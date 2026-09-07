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

package monitors

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// The SDK adds a zero-valued element to a TypeSet that has a diff - https://github.com/hashicorp/terraform-plugin-sdk/issues/895.
// Such an element carries either no values at all or nothing but the schema defaults, but never a key, so it must not reach the API.
func TestTagsWithSourceInfo_UnmarshalHCL_DropsPhantomEmptyTag(t *testing.T) {
	cases := []struct {
		name    string
		phantom map[string]any
	}{
		{
			name:    "without values",
			phantom: map[string]any{"source": "", "context": "", "key": "", "value": ""},
		},
		{
			name:    "with the schema defaults",
			phantom: map[string]any{"source": "USER", "context": "CONTEXTLESS", "key": "", "value": ""},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, new(TagsWithSourceInfo).Schema(), map[string]any{
				"tag": []any{
					map[string]any{"source": "USER", "context": "CONTEXTLESS", "key": "protocol", "value": "TCP"},
					c.phantom,
				},
			})

			var tags TagsWithSourceInfo
			err := tags.UnmarshalHCL(hcl.DecoderFrom(d))

			assert.NoError(t, err)
			assert.Equal(t, TagsWithSourceInfo{
				{Source: &defaultTagSource, Context: &defaultTagContext, Key: "protocol", Value: new("TCP")},
			}, tags)
		})
	}
}

// The API assigns `source` and `context` to every tag it stores, so a configuration that omits them
// has to plan to those very values - otherwise every refresh reports drift and every plan reverts it.
func TestTagsWithSourceInfo_NoPlanDiffWhenSourceAndContextAreOmitted(t *testing.T) {
	tagsResource := &schema.Resource{Schema: map[string]*schema.Schema{"tags": new(Settings).Schema()["tags"]}}

	stateTag := func(key string, value string) map[string]any {
		return map[string]any{"source": "USER", "context": "CONTEXTLESS", "key": key, "value": value}
	}
	data := tagsResource.TestResourceData()
	data.SetId("MULTI_PROTOCOL_TEST-1234567890000000")
	assert.NoError(t, data.Set("tags", []any{map[string]any{"tag": []any{
		stateTag("protocol", "TCP"),
		stateTag("test-decentralize", ""),
	}}}))

	configuredTag := func(key string, value string) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{
			"source":  cty.NullVal(cty.String),
			"context": cty.NullVal(cty.String),
			"key":     cty.StringVal(key),
			"value":   cty.StringVal(value),
		})
	}
	configuration := terraform.NewResourceConfigShimmed(cty.ObjectVal(map[string]cty.Value{
		"id": cty.NullVal(cty.String),
		"tags": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"tag": cty.SetVal([]cty.Value{
				configuredTag("protocol", "TCP"),
				configuredTag("test-decentralize", ""),
			}),
		})}),
	}), tagsResource.CoreConfigSchema())

	diff, err := tagsResource.SimpleDiff(t.Context(), data.State(), configuration, nil)

	assert.NoError(t, err)
	assert.Empty(t, diff.Attributes)
}
