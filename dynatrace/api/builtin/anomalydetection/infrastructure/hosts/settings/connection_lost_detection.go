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

package hosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConnectionLostDetection struct {
	Enabled             bool                                `json:"enabled"`                       // Detect host or monitoring connection lost problems
	OnGracefulShutdowns *ConnectionLostDetectionSensitivity `json:"onGracefulShutdowns,omitempty"` // Graceful host shutdowns
}

func (me *ConnectionLostDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Detect host or monitoring connection lost problems",
			Required:    true,
		},
		"on_graceful_shutdowns": {
			Type:        schema.TypeString,
			Description: "Graceful host shutdowns. Possible values: `DONT_ALERT_ON_GRACEFUL_SHUTDOWN`, `ALERT_ON_GRACEFUL_SHUTDOWN`",
			Optional:    true,
		},
	}
}

func (me *ConnectionLostDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":               me.Enabled,
		"on_graceful_shutdowns": me.OnGracefulShutdowns,
	})
}

func (me *ConnectionLostDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":               &me.Enabled,
		"on_graceful_shutdowns": &me.OnGracefulShutdowns,
	})
}
