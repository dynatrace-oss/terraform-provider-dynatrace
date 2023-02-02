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

package databases

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DatabaseConnections struct {
	Enabled           bool `json:"enabled"`                     // Detect failed database connects
	MaxFailedConnects *int `json:"maxFailedConnects,omitempty"` // Threshold
	TimePeriod        *int `json:"timePeriod,omitempty"`        // Time span
}

func (me *DatabaseConnections) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Detect failed database connects",
			Required:    true,
		},
		"max_failed_connects": {
			Type:        schema.TypeInt,
			Description: "Threshold",
			Optional:    true,
		},
		"time_period": {
			Type:        schema.TypeInt,
			Description: "Time span",
			Optional:    true,
		},
	}
}

func (me *DatabaseConnections) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":             me.Enabled,
		"max_failed_connects": me.MaxFailedConnects,
		"time_period":         me.TimePeriod,
	})
}

func (me *DatabaseConnections) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":             &me.Enabled,
		"max_failed_connects": &me.MaxFailedConnects,
		"time_period":         &me.TimePeriod,
	})
}
