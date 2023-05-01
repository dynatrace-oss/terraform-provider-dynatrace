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

type Settings struct {
	ApiColor      string      `json:"apiColor"`             // This color will be used to highlight APIs when viewing code level data, such as distributed traces or method hotspots.
	ApiName       string      `json:"apiName"`              // API name
	Conditions    ApiRules    `json:"conditions,omitempty"` // List of conditions
	Technology    *Technology `json:"technology,omitempty"` // Restrict this rule to a specific technology.
	ThirdPartyApi bool        `json:"thirdPartyApi"`        // This API defines a third party library
}

func (me *Settings) Name() string {
	return me.ApiName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_color": {
			Type:        schema.TypeString,
			Description: "This color will be used to highlight APIs when viewing code level data, such as distributed traces or method hotspots.",
			Required:    true,
		},
		"api_name": {
			Type:        schema.TypeString,
			Description: "API name",
			Required:    true,
		},
		"conditions": {
			Type:        schema.TypeList,
			Description: "List of conditions",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(ApiRules).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"technology": {
			Type:        schema.TypeString,
			Description: "Restrict this rule to a specific technology.",
			Optional:    true,
		},
		"third_party_api": {
			Type:        schema.TypeBool,
			Description: "This API defines a third party library",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"api_color":       me.ApiColor,
		"api_name":        me.ApiName,
		"conditions":      me.Conditions,
		"technology":      me.Technology,
		"third_party_api": me.ThirdPartyApi,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"api_color":       &me.ApiColor,
		"api_name":        &me.ApiName,
		"conditions":      &me.Conditions,
		"technology":      &me.Technology,
		"third_party_api": &me.ThirdPartyApi,
	})
}
