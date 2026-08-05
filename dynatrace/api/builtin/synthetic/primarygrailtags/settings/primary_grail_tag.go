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

package primarygrailtags

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PrimaryGrailTags []*PrimaryGrailTag

func (me *PrimaryGrailTags) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tag": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(PrimaryGrailTag).Schema()},
		},
	}
}

func (me PrimaryGrailTags) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("tag", me)
}

func (me *PrimaryGrailTags) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("tag", me)
}

type PrimaryGrailTag struct {
	Key   string `json:"key"`   // Primary grail tag's key
	Value string `json:"value"` // Primary grail tag's value
}

func (me *PrimaryGrailTag) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "Primary grail tag's key",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Primary grail tag's value",
			Required:    true,
		},
	}
}

func (me *PrimaryGrailTag) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":   me.Key,
		"value": me.Value,
	})
}

func (me *PrimaryGrailTag) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":   &me.Key,
		"value": &me.Value,
	})
}
