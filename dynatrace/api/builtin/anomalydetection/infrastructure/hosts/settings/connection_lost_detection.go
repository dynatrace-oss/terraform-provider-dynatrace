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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConnectionLostDetection struct {
	Enabled             bool                                `json:"enabled"`                       // This setting is enabled (`true`) or disabled (`false`)
	OnGracefulShutdowns *ConnectionLostDetectionSensitivity `json:"onGracefulShutdowns,omitempty"` // Graceful host shutdowns. Possible Values: `ALERT_ON_GRACEFUL_SHUTDOWN`, `DONT_ALERT_ON_GRACEFUL_SHUTDOWN`
}

func (me *ConnectionLostDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"on_graceful_shutdowns": {
			Type:        schema.TypeString,
			Description: "Graceful host shutdowns. Possible Values: `ALERT_ON_GRACEFUL_SHUTDOWN`, `DONT_ALERT_ON_GRACEFUL_SHUTDOWN`",
			Optional:    true, // precondition
		},
	}
}

func (me *ConnectionLostDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":               me.Enabled,
		"on_graceful_shutdowns": me.OnGracefulShutdowns,
	})
}

func (me *ConnectionLostDetection) HandlePreconditions() error {
	if (me.OnGracefulShutdowns == nil) && (me.Enabled) {
		return fmt.Errorf("'on_graceful_shutdowns' must be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if (me.OnGracefulShutdowns != nil) && (!me.Enabled) {
		return fmt.Errorf("'on_graceful_shutdowns' must not be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	return nil
}

func (me *ConnectionLostDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":               &me.Enabled,
		"on_graceful_shutdowns": &me.OnGracefulShutdowns,
	})
}
