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

package nodes

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ID                     string   `json:"entityId"`               // The ID of a node
	Hostname               string   `json:"hostname"`               // The hostname of a node
	IPs                    []string `json:"ips"`                    // The known IP addresses of the node
	Version                string   `json:"version"`                // The version of a node
	BrowserMonitorsEnabled bool     `json:"browserMonitorsEnabled"` // Browser Monitors are enabled (`true`) or not (`false`)
	ActiveGateVersion      string   `json:"activeGateVersion"`      // The version of the Active Gate
	OneAgentRoutingEnabled bool     `json:"oneAgentRoutingEnabled"` // The Active Gate has the One Agent routing enabled (`true`) or not (`false`)
	OperatingSystem        string   `json:"operatingSystem"`        // The Active Gate's host operating system
	AutoUpdateEnabled      bool     `json:"autoUpdateEnabled"`      // The Active Gate has the Auto update option enabled (`true`) or not (`false`)
	Status                 string   `json:"status"`                 // The status of the synthetic node
	PlayerVersion          string   `json:"playerVersion"`          // The version of the synthetic player
	HealthCheckStatus      string   `json:"healthCheckStatus"`      // The health check status of the synthetic node
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of a node for usage within `dynatrace_synthetic_location`",
			Computed:    true,
		},
		"hostname": {
			Type:        schema.TypeString,
			Description: "The hostname of a node",
			Computed:    true,
		},
		"ips": {
			Type:        schema.TypeSet,
			Description: "The known IP addresses of the node",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
		},
		"version": {
			Type:        schema.TypeString,
			Description: "The version of a node",
			Computed:    true,
		},
		"browser_monitors": {
			Type:        schema.TypeBool,
			Description: "Specifies whether Browser Monitors are enabled or not",
			Computed:    true,
		},
		"active_gate_version": {
			Type:        schema.TypeString,
			Description: "The version of the Active Gate",
			Computed:    true,
		},
		"one_agent_routing": {
			Type:        schema.TypeBool,
			Description: "Specifies whether the Active Gate has the One Agent routing enabled",
			Computed:    true,
		},
		"operating_system": {
			Type:        schema.TypeString,
			Description: "The Active Gate's host operating system",
			Computed:    true,
		},
		"auto_update": {
			Type:        schema.TypeBool,
			Description: "Specifies whether the Active Gate  has the Auto update option enabled",
			Computed:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "The status of the synthetic node",
			Computed:    true,
		},
		"player_version": {
			Type:        schema.TypeString,
			Description: "The version of the synthetic player",
			Computed:    true,
		},
		"health_check_status": {
			Type:        schema.TypeString,
			Description: "The health check status of the synthetic node",
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":                  me.ID,
		"hostname":            me.Hostname,
		"ips":                 me.IPs,
		"version":             me.Version,
		"browser_monitors":    me.BrowserMonitorsEnabled,
		"active_gate_version": me.ActiveGateVersion,
		"one_agent_routing":   me.OneAgentRoutingEnabled,
		"operating_system":    me.OperatingSystem,
		"auto_update":         me.AutoUpdateEnabled,
		"status":              me.Status,
		"player_version":      me.PlayerVersion,
		"health_check_status": me.HealthCheckStatus,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":                  &me.ID,
		"hostname":            &me.Hostname,
		"ips":                 &me.IPs,
		"version":             &me.Version,
		"browser_monitors":    &me.BrowserMonitorsEnabled,
		"active_gate_version": &me.ActiveGateVersion,
		"one_agent_routing":   &me.OneAgentRoutingEnabled,
		"operating_system":    &me.OperatingSystem,
		"auto_update":         &me.AutoUpdateEnabled,
		"status":              &me.Status,
		"player_version":      &me.PlayerVersion,
		"health_check_status": &me.HealthCheckStatus,
	})
}
