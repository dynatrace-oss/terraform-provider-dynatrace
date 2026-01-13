/**
* @license
* Copyright 2025 Dynatrace LLC
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

package extension

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled bool   `json:"diskDeviceMonitoringExtensionActive"` // The Disk Analytics feature requires an extension to be added to your environment. You can add the Disk Analytics extension to your environment from [Dynatrace Hub](/ui/hub/ext/com.dynatrace.extension.disk-devices#information). The Disk Analytics extension consumes custom metrics and [Davis data units](https://www.dynatrace.com/support/help/shortlink/metric-cost-calculation).\n\nAfter you have added the Disk Analytics extension, you can enable the Data Collection in host or host-group level settings. If you enable the Data Collection without adding the extension the data is only visible in the data explorer.\n\nFor details, see [Disk Analytics extension documentation](https://dt-url.net/3a03v9v).
	Scope   string `json:"-" scope:"scope"`                     // The scope of this setting (HOST, HOST_GROUP)
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The Disk Analytics feature requires an extension to be added to your environment. You can add the Disk Analytics extension to your environment from [Dynatrace Hub](/ui/hub/ext/com.dynatrace.extension.disk-devices#information). The Disk Analytics extension consumes custom metrics and [Davis data units](https://www.dynatrace.com/support/help/shortlink/metric-cost-calculation).\n\nAfter you have added the Disk Analytics extension, you can enable the Data Collection in host or host-group level settings. If you enable the Data Collection without adding the extension the data is only visible in the data explorer.\n\nFor details, see [Disk Analytics extension documentation](https://dt-url.net/3a03v9v).",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP)",
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"scope":   me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled": &me.Enabled,
		"scope":   &me.Scope,
	})
}
