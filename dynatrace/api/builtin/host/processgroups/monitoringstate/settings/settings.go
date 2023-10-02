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

package monitoringstate

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	HostID          string                     `json:"-" scope:"hostId"` // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	MonitoringState ProcessGroupMonitoringMode `json:"MonitoringState"`  // Possible Values: `DEFAULT`, `MONITORING_OFF`, `MONITORING_ON`
	ProcessGroup    string                     `json:"ProcessGroup"`     // Process group
}

func (me *Settings) Name() string {
	return me.HostID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
			ForceNew:    true,
		},
		"monitoring_state": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DEFAULT`, `MONITORING_OFF`, `MONITORING_ON`",
			Required:    true,
		},
		"process_group": {
			Type:        schema.TypeString,
			Description: "Process group",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"host_id":          me.HostID,
		"monitoring_state": me.MonitoringState,
		"process_group":    me.ProcessGroup,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"host_id":          &me.HostID,
		"monitoring_state": &me.MonitoringState,
		"process_group":    &me.ProcessGroup,
	})
}
