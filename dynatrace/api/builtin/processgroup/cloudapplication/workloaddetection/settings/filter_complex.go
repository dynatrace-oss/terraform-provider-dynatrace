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

package workloaddetection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FilterComplexes []*FilterComplex

func (me *FilterComplexes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(FilterComplex).Schema()},
		},
	}
}

func (me FilterComplexes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *FilterComplexes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter", me)
}

type FilterComplex struct {
	Enabled          bool              `json:"enabled"`          // This setting is enabled (`true`) or disabled (`false`)
	InclusionToggles *InclusionToggles `json:"inclusionToggles"` // ID calculation based on
	MatchFilter      *MatchFilter      `json:"matchFilter"`      // When namespace
}

func (me *FilterComplex) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"inclusion_toggles": {
			Type:        schema.TypeList,
			Description: "ID calculation based on",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(InclusionToggles).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"match_filter": {
			Type:        schema.TypeList,
			Description: "When namespace",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(MatchFilter).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *FilterComplex) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":           me.Enabled,
		"inclusion_toggles": me.InclusionToggles,
		"match_filter":      me.MatchFilter,
	})
}

func (me *FilterComplex) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":           &me.Enabled,
		"inclusion_toggles": &me.InclusionToggles,
		"match_filter":      &me.MatchFilter,
	})
}
