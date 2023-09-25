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

package rulesettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled  bool              `json:"enabled"`  // This setting is enabled (`true`) or disabled (`false`)
	Mode     MonitoringMode    `json:"mode"`     // Possible Values: `MONITORING_OFF`, `MONITORING_ON`
	Operator ConditionOperator `json:"operator"` // Possible Values: `EQUALS`, `NOT_EQUALS`
	Property Property          `json:"property"` // Possible Values: `HOST_TAG`, `MANAGEMENT_ZONE`, `PROCESS_TAG`
	Value    string            `json:"value"`    // Condition value
}

func (me *Settings) Name() string {
	return me.Value
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `MONITORING_OFF`, `MONITORING_ON`",
			Required:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Possible Values: `EQUALS`, `NOT_EQUALS`",
			Required:    true,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Possible Values: `HOST_TAG`, `MANAGEMENT_ZONE`, `PROCESS_TAG`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Condition value",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":  me.Enabled,
		"mode":     me.Mode,
		"operator": me.Operator,
		"property": me.Property,
		"value":    me.Value,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":  &me.Enabled,
		"mode":     &me.Mode,
		"operator": &me.Operator,
		"property": &me.Property,
		"value":    &me.Value,
	})
}
