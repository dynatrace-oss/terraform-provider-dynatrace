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

package privacy

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataCollection struct {
	OptInModeEnabled bool `json:"optInModeEnabled"` // With [Data-collection and opt-in mode](https://dt-url.net/7l3p0p3h) enabled, Real User Monitoring data isn't captured until dtrum.enable() is called for specific user sessions.
}

func (me *DataCollection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"opt_in_mode_enabled": {
			Type:        schema.TypeBool,
			Description: "With [Data-collection and opt-in mode](https://dt-url.net/7l3p0p3h) enabled, Real User Monitoring data isn't captured until dtrum.enable() is called for specific user sessions.",
			Required:    true,
		},
	}
}

func (me *DataCollection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"opt_in_mode_enabled": me.OptInModeEnabled,
	})
}

func (me *DataCollection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"opt_in_mode_enabled": &me.OptInModeEnabled,
	})
}
