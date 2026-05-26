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

package rumweb

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AppTrafficDrops struct {
	Enabled      bool          `json:"enabled"`                // This setting is enabled (`true`) or disabled (`false`)
	TrafficDrops *TrafficDrops `json:"trafficDrops,omitempty"` // Dynatrace learns your typical application traffic over an observation period of one week.\n\n  Depending on this expected value Dynatrace detects abnormal traffic drops within your application.
}

func (me *AppTrafficDrops) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"traffic_drops": {
			Type:        schema.TypeList,
			Description: "Dynatrace learns your typical application traffic over an observation period of one week.\n\n  Depending on this expected value Dynatrace detects abnormal traffic drops within your application.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(TrafficDrops).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *AppTrafficDrops) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":       me.Enabled,
		"traffic_drops": me.TrafficDrops,
	})
}

func (me *AppTrafficDrops) HandlePreconditions() error {
	if (me.TrafficDrops != nil) && (!me.Enabled) {
		return fmt.Errorf("'traffic_drops' must not be specified unless 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.TrafficDrops == nil) && (me.Enabled) {
		return fmt.Errorf("'traffic_drops' must be specified when 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	return nil
}

func (me *AppTrafficDrops) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":       &me.Enabled,
		"traffic_drops": &me.TrafficDrops,
	})
}
