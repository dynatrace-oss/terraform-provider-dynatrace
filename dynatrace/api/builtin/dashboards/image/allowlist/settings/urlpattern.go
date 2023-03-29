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

type URLPatterns []*URLPattern

func (me *URLPatterns) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"urlpattern": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(URLPattern).Schema()},
		},
	}
}

func (me URLPatterns) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("urlpattern", me)
}

func (me *URLPatterns) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("urlpattern", me)
}

type URLPattern struct {
	Rule     RuleEnum `json:"rule"`     // Possible Values: `Equals`, `StartsWith`
	Template string   `json:"template"` // Pattern
}

func (me *URLPattern) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Equals`, `StartsWith`",
			Required:    true,
		},
		"template": {
			Type:        schema.TypeString,
			Description: "Pattern",
			Required:    true,
		},
	}
}

func (me *URLPattern) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"rule":     me.Rule,
		"template": me.Template,
	})
}

func (me *URLPattern) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"rule":     &me.Rule,
		"template": &me.Template,
	})
}
