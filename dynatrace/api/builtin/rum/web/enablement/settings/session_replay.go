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

package enablement

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SessionReplay struct {
	CostAndTrafficControl int  `json:"costAndTrafficControl"` // (Field has overlap with `dynatrace_web_application`) [Percentage of user sessions recorded with Session Replay](https://dt-url.net/sr-cost-traffic-control). For example, if you have 50% for RUM and 50% for Session Replay, it results in 25% of sessions recorded with Session Replay.
	Enabled               bool `json:"enabled"`               // (Field has overlap with `dynatrace_web_application`) This setting is enabled (`true`) or disabled (`false`)
}

func (me *SessionReplay) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cost_and_traffic_control": {
			Type:        schema.TypeInt,
			Description: "(Field has overlap with `dynatrace_web_application`) [Percentage of user sessions recorded with Session Replay](https://dt-url.net/sr-cost-traffic-control). For example, if you have 50% for RUM and 50% for Session Replay, it results in 25% of sessions recorded with Session Replay.",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "(Field has overlap with `dynatrace_web_application`) This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
	}
}

func (me *SessionReplay) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cost_and_traffic_control": me.CostAndTrafficControl,
		"enabled":                  me.Enabled,
	})
}

func (me *SessionReplay) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cost_and_traffic_control": &me.CostAndTrafficControl,
		"enabled":                  &me.Enabled,
	})
}
