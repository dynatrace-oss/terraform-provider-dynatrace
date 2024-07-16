/**
* @license
* Copyright 2020 Dynatrace LLC
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TagsWithSourceInfo []*TagWithSourceInfo

func (me *TagsWithSourceInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tag": {
			Type:        schema.TypeSet,
			Description: "Tag with source of a Dynatrace entity.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TagWithSourceInfo).Schema()},
		},
	}
}

func (me TagsWithSourceInfo) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("tag", me)
}

func (me *TagsWithSourceInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("tag", me)
}

// Tag with source of a Dynatrace entity
type TagWithSourceInfo struct {
	Source  *TagSource `json:"source,omitempty"`  // The source of the tag, possible values: `AUTO`, `RULE_BASED` or `USER`
	Context *string    `json:"context,omitempty"` // The origin of the tag, such as AWS or Cloud Foundry.\n\nCustom tags use the CONTEXTLESS value
	Key     string     `json:"key"`               // The key of the tag
	Value   *string    `json:"value,omitempty"`   // The value of the tag
}

func (me *TagWithSourceInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": {
			Type:        schema.TypeString,
			Description: "The source of the tag, possible values: `AUTO`, `RULE_BASED` or `USER`",
			Optional:    true,
		},
		"context": {
			Type:        schema.TypeString,
			Description: "The origin of the tag, such as AWS or Cloud Foundry.\n\nCustom tags use the CONTEXTLESS value",
			Optional:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The key of the tag",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: " The value of the tag",
			Optional:    true,
		},
	}
}

func (me TagWithSourceInfo) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"source":  me.Source,
		"context": me.Context,
		"key":     me.Key,
		"value":   me.Value,
	})
}

func (me *TagWithSourceInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"source":  &me.Source,
		"context": &me.Context,
		"key":     &me.Key,
		"value":   &me.Value,
	})
}
