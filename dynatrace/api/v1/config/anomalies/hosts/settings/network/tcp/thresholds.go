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

package tcp

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NetworkTcpProblemsThresholds Custom thresholds for TCP connection problems. If not set, automatic mode is used.
//
//	**All** of these conditions must be met to trigger an alert.
type Thresholds struct {
	NewConnectionFailuresPercentage  int32 `json:"newConnectionFailuresPercentage"`  // Percentage of new connection failures is higher than *X*% in 3 out of 5 samples.
	FailedConnectionsNumberPerMinute int32 `json:"failedConnectionsNumberPerMinute"` // Number of failed connections is higher than *X* connections per minute in 3 out of 5 samples.
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"new_connection_failures": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Percentage of new connection failures is higher than *X*% in 3 out of 5 samples",
		},
		"failed_connections": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of failed connections is higher than *X* connections per minute in 3 out of 5 samples",
		},
	}
}

func (me *Thresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"new_connection_failures": me.NewConnectionFailuresPercentage,
		"failed_connections":      me.FailedConnectionsNumberPerMinute,
	})
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("new_connection_failures"); ok {
		me.NewConnectionFailuresPercentage = int32(value.(int))
	}
	if value, ok := decoder.GetOk("failed_connections"); ok {
		me.FailedConnectionsNumberPerMinute = int32(value.(int))
	}
	return nil
}
