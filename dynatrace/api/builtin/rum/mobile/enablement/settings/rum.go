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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rum struct {
	CostAndTrafficControl int   `json:"costAndTrafficControl"`    // Percentage of user sessions captured and analyzed. By default, Dynatrace captures all user actions and user sessions for analysis. This approach ensures complete insight into your application’s performance and customer experience. You can optionally reduce the granularity of user-action and user-session analysis by capturing a lower percentage of user sessions. While this approach can reduce monitoring costs, it also results in lower visibility into how your customers are using your applications. For example, a setting of 10% results in Dynatrace analyzing only every tenth user session.
	Enabled               bool  `json:"enabled"`                  // This setting is enabled (`true`) or disabled (`false`)
	EnabledOnGrail        *bool `json:"enabledOnGrail,omitempty"` // Please be aware that only mobile agents with version **8.303 or higher** can ingest Grail events
}

func (me *Rum) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cost_and_traffic_control": {
			Type:        schema.TypeInt,
			Description: "Percentage of user sessions captured and analyzed. By default, Dynatrace captures all user actions and user sessions for analysis. This approach ensures complete insight into your application’s performance and customer experience. You can optionally reduce the granularity of user-action and user-session analysis by capturing a lower percentage of user sessions. While this approach can reduce monitoring costs, it also results in lower visibility into how your customers are using your applications. For example, a setting of 10% results in Dynatrace analyzing only every tenth user session.",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"enabled_on_grail": {
			Type:        schema.TypeBool,
			Description: "Please be aware that only mobile agents with version **8.303 or higher** can ingest Grail events",
			Optional:    true, // nullable & precondition
		},
	}
}

func (me *Rum) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cost_and_traffic_control": me.CostAndTrafficControl,
		"enabled":                  me.Enabled,
		"enabled_on_grail":         me.EnabledOnGrail,
	})
}

func (me *Rum) HandlePreconditions() error {
	if (me.EnabledOnGrail == nil) && (me.Enabled) {
		me.EnabledOnGrail = opt.NewBool(false)
	}
	return nil
}

func (me *Rum) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cost_and_traffic_control": &me.CostAndTrafficControl,
		"enabled":                  &me.Enabled,
		"enabled_on_grail":         &me.EnabledOnGrail,
	})
}
