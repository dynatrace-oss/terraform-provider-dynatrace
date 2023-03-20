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

package monitoringrule

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled  bool              `json:"enabled"`         // This setting is enabled (`true`) or disabled (`false`)
	Mode     MonitoringMode    `json:"mode"`            // Possible Values: `MONITORING_OFF`, `MONITORING_ON`
	Operator ConditionOperator `json:"operator"`        // Possible Values: `CONTAINS`, `ENDS`, `EQUALS`, `EXISTS`, `NOT_CONTAINS`, `NOT_ENDS`, `NOT_EQUALS`, `NOT_EXISTS`, `NOT_STARTS`, `STARTS`
	Property ContainerItem     `json:"property"`        // Possible Values: `CONTAINER_NAME`, `IMAGE_NAME`, `KUBERNETES_BASEPODNAME`, `KUBERNETES_CONTAINERNAME`, `KUBERNETES_FULLPODNAME`, `KUBERNETES_NAMESPACE`, `KUBERNETES_PODUID`
	Value    *string           `json:"value,omitempty"` // Condition value
}

func (me *Settings) Name() string {
	return "container_rule"
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
			Description: "Possible Values: `CONTAINS`, `ENDS`, `EQUALS`, `EXISTS`, `NOT_CONTAINS`, `NOT_ENDS`, `NOT_EQUALS`, `NOT_EXISTS`, `NOT_STARTS`, `STARTS`",
			Required:    true,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINER_NAME`, `IMAGE_NAME`, `KUBERNETES_BASEPODNAME`, `KUBERNETES_CONTAINERNAME`, `KUBERNETES_FULLPODNAME`, `KUBERNETES_NAMESPACE`, `KUBERNETES_PODUID`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Condition value",
			Optional:    true,
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
