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

package allowlist

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AllowRules []*AllowRule

func (me *AllowRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_rule": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(AllowRule).Schema()},
		},
	}
}

func (me AllowRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("allow_rule", me)
}

func (me *AllowRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("allow_rule", me)
}

type AllowRule struct {
	Enabled   bool    `json:"enabled"`   // This setting is enabled (`true`) or disabled (`false`)
	MetricKey string  `json:"metricKey"` // Metric key
	Pattern   Pattern `json:"pattern"`   // Possible Values: `CONTAINS`, `EQUALS`, `STARTSWITH`
}

func (me *AllowRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"metric_key": {
			Type:        schema.TypeString,
			Description: "Metric key",
			Required:    true,
		},
		"pattern": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `EQUALS`, `STARTSWITH`",
			Required:    true,
		},
	}
}

func (me *AllowRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":    me.Enabled,
		"metric_key": me.MetricKey,
		"pattern":    me.Pattern,
	})
}

func (me *AllowRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":    &me.Enabled,
		"metric_key": &me.MetricKey,
		"pattern":    &me.Pattern,
	})
}
