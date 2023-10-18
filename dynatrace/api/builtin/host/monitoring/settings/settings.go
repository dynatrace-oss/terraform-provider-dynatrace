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

package monitoring

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AutoInjection bool   `json:"autoInjection"`    // An auto-injection disabled with [oneagentctl](https://dt-url.net/oneagentctl) takes precedence over this setting and cannot be changed from the Dynatrace web UI.
	Enabled       bool   `json:"enabled"`          // This setting is enabled (`true`) or disabled (`false`)
	HostID        string `json:"-" scope:"hostId"` // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
}

func (me *Settings) Name() string {
	return me.HostID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_injection": {
			Type:        schema.TypeBool,
			Description: "An auto-injection disabled with [oneagentctl](https://dt-url.net/oneagentctl) takes precedence over this setting and cannot be changed from the Dynatrace web UI.",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"full_stack": {
			Type:        schema.TypeBool,
			Description: "Dynatrace uses full-stack monitoring by default, to monitor every aspect of your environment, including all processes, services, and applications detected on your hosts. \n\nIf you turn off full-stack monitoring, Dynatrace will only monitor your infrastructure. You will lose access to application performance, user experience data, code-level visibility and PurePath insights. \n\nTo learn more, visit [Infrastructure Monitoring mode](https://www.dynatrace.com/support/help/shortlink/infrastructure).\n\nPlease note that changing the monitoring mode will impact the license consumption of this OneAgent. To learn more, visit [Host units](https://dt-url.net/hi03uns).",
			Optional:    true,
			Deprecated:  "This attribute is not supported anymore by the Dynatrace API",
		},
		"host_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_injection": me.AutoInjection,
		"enabled":        me.Enabled,
		"host_id":        me.HostID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_injection": &me.AutoInjection,
		"enabled":        &me.Enabled,
		"host_id":        &me.HostID,
	})
}
