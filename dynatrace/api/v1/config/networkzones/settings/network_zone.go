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

package networkzones

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NetworkZone TODO: documentation
type NetworkZone struct {
	ID                           *string  `json:"id,omitempty"`                           // The ID of the network zone
	Description                  *string  `json:"description,omitempty"`                  // A short description of the network zone
	AltZones                     []string `json:"alternativeZones,omitempty"`             // A list of alternative network zones.
	NumOfOneAgentsFromOtherZones *int     `json:"numOfOneAgentsFromOtherZones,omitempty"` // The number of OneAgents from other network zones that are using ActiveGates in the network zone.
	NumOfOneAgentsUsing          *int     `json:"numOfOneAgentsUsing,omitempty"`          // The number of OneAgents that are using ActiveGates in the network zone.
	NumofConfiguredActiveGates   *int     `json:"numOfConfiguredActiveGates,omitempty"`   // The number of ActiveGates in the network zone.
	NumOfConfiguredOneAgents     *int     `json:"numOfConfiguredOneAgents,omitempty"`     // The number of OneAgents that are configured to use the network zone as primary.
}

func (me *NetworkZone) Name() string {
	return *me.ID
}

type NetworkZones struct {
	Zones []NetworkZone `json:"networkZones"`
}

func (me *NetworkZone) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the network zone",
			Optional:    true,
		},
		"alternative_zones": {
			Type:        schema.TypeList,
			Description: "A list of alternative network zones.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"num_of_oneagents_from_other_zones": {
			Type:        schema.TypeInt,
			Description: "The number of OneAgents from other network zones that are using ActiveGates in the network zone.",
			Optional:    true,
		},
		"num_of_oneagents_using": {
			Type:        schema.TypeInt,
			Description: "The number of OneAgents that are using ActiveGates in the network zone.",
			Optional:    true,
		},
		"num_of_configured_activegates": {
			Type:        schema.TypeInt,
			Description: "The number of ActiveGates in the network zone.",
			Optional:    true,
		},
		"num_of_configured_oneagents": {
			Type:        schema.TypeInt,
			Description: "The number of OneAgents that are configured to use the network zone as primary.",
			Optional:    true,
		},
	}
}

func (me *NetworkZone) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":                       me.Description,
		"alternative_zones":                 me.AltZones,
		"num_of_oneagents_from_other_zones": me.NumOfOneAgentsFromOtherZones,
		"num_of_oneagents_using":            me.NumOfOneAgentsUsing,
		"num_of_configured_activegates":     me.NumofConfiguredActiveGates,
		"num_of_configured_oneagents":       me.NumOfConfiguredOneAgents,
	})
}

func (me *NetworkZone) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description":       &me.Description,
		"alternative_zones": &me.AltZones,
	})
}
