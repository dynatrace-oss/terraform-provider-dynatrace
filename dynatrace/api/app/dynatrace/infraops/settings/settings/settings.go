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

package infraopssettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Show_monitoring_candidates bool `json:"show.monitoring.candidates"` // When set to true, the app will display monitoring candidates in the Hosts table
	Show_standalone_hosts      bool `json:"show.standalone.hosts"`      // When set to true, the app will display app only hosts in the Hosts table
}

func (me *Settings) Name() string {
	return "environment"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"show_monitoring_candidates": {
			Type:        schema.TypeBool,
			Description: "When set to true, the app will display monitoring candidates in the Hosts table",
			Required:    true,
		},
		"show_standalone_hosts": {
			Type:        schema.TypeBool,
			Description: "When set to true, the app will display app only hosts in the Hosts table",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"show_monitoring_candidates": me.Show_monitoring_candidates,
		"show_standalone_hosts":      me.Show_standalone_hosts,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"show_monitoring_candidates": &me.Show_monitoring_candidates,
		"show_standalone_hosts":      &me.Show_standalone_hosts,
	})
}
