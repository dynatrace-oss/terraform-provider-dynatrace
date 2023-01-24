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

package detection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Rule the configuration of an application detection rule
type Rule struct {
	Name                  *string       `json:"name,omitempty"`        // the unique name of the Application detection rule
	Order                 *string       `json:"order,omitempty"`       // the order of the rule in the rules list
	ApplicationIdentifier string        `json:"applicationIdentifier"` // the Dynatrace entity ID of the application, for example APPLICATION-4A3B43
	FilterConfig          *FilterConfig `json:"filterConfig"`          // the condition of an application detection rule
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The unique name of the Application detection rule",
		},
		"order": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The order of the rule in the rules list",
		},
		"application_identifier": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Dynatrace entity ID of the application, for example APPLICATION-4A3B43",
		},
		"filter_config": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FilterConfig).Schema()},
			Description: "The condition of an application detection rule",
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":                   me.Name,
		"order":                  me.Order,
		"application_identifier": me.ApplicationIdentifier,
		"filter_config":          me.FilterConfig,
	})
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":                   &me.Name,
		"order":                  &me.Order,
		"application_identifier": &me.ApplicationIdentifier,
		"filter_config":          &me.FilterConfig,
	})
}
