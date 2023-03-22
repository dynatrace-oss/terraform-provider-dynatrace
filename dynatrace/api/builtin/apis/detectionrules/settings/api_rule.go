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

package detectionrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ApiRules []*ApiRule

func (me *ApiRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ApiRule).Schema()},
		},
	}
}

func (me ApiRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("condition", me)
}

func (me *ApiRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("condition", me)
}

type ApiRule struct {
	Base    Base    `json:"base"`    // Possible Values: `FILE_NAME`, `FQCN`, `PACKAGE`
	Matcher Matcher `json:"matcher"` // Possible Values: `BEGINS_WITH`, `CONTAINS`
	Pattern string  `json:"pattern"`
}

func (me *ApiRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"base": {
			Type:        schema.TypeString,
			Description: "Possible Values: `FILE_NAME`, `FQCN`, `PACKAGE`",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BEGINS_WITH`, `CONTAINS`",
			Required:    true,
		},
		"pattern": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *ApiRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"base":    me.Base,
		"matcher": me.Matcher,
		"pattern": me.Pattern,
	})
}

func (me *ApiRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"base":    &me.Base,
		"matcher": &me.Matcher,
		"pattern": &me.Pattern,
	})
}
