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

package externalwebservice

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Conditions      Conditions          `json:"conditions,omitempty"`      // A list of conditions necessary for the rule to take effect. If multiple conditions are specified, they must **all** match a Request for the rule to apply. Conditions are evaluated against attributes, but do not modify them.
	Description     *string             `json:"description,omitempty"`     // Description
	Enabled         bool                `json:"enabled"`                   // This setting is enabled (`true`) or disabled (`false`)
	IdContributors  *IdContributorsType `json:"idContributors"`            // Contributors to the Service Identifier calculation. URL path is always applied as an Id Contributor. You can exclude the port contribution by disabling the switch.
	ManagementZones []string            `json:"managementZones,omitempty"` // Define a management zone filter for this service detection rule.
	Name            string              `json:"name"`                      // Rule name
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"conditions": {
			Type:        schema.TypeList,
			Description: "A list of conditions necessary for the rule to take effect. If multiple conditions are specified, they must **all** match a Request for the rule to apply. Conditions are evaluated against attributes, but do not modify them.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Conditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"id_contributors": {
			Type:        schema.TypeList,
			Description: "Contributors to the Service Identifier calculation. URL path is always applied as an Id Contributor. You can exclude the port contribution by disabling the switch.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(IdContributorsType).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"management_zones": {
			Type:        schema.TypeSet,
			Description: "Define a management zone filter for this service detection rule.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"conditions":       me.Conditions,
		"description":      me.Description,
		"enabled":          me.Enabled,
		"id_contributors":  me.IdContributors,
		"management_zones": me.ManagementZones,
		"name":             me.Name,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"conditions":       &me.Conditions,
		"description":      &me.Description,
		"enabled":          &me.Enabled,
		"id_contributors":  &me.IdContributors,
		"management_zones": &me.ManagementZones,
		"name":             &me.Name,
	})
}
