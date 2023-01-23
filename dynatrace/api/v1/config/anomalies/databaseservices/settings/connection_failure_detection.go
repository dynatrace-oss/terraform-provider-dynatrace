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

package databaseservices

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ConnectionFailureDetection Parameters of the failed database connections detection.
// The alert is triggered when failed connections number exceeds **connectionFailsCount** during any **timePeriodMinutes** minutes period.
type ConnectionFailureDetection struct {
	ConnectionFailsCount *int32 `json:"connectionFailsCount,omitempty"` // Number of failed database connections during any **timePeriodMinutes** minutes period to trigger an alert.
	Enabled              bool   `json:"enabled"`                        // The detection is enabled (`true`) or disabled (`false`).
	TimePeriodMinutes    *int32 `json:"timePeriodMinutes,omitempty"`    // The *X* minutes time period during which the **connectionFailsCount** is evaluated.
}

func (me *ConnectionFailureDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_fails_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of failed database connections during any **eval_period** minutes period to trigger an alert",
		},
		"eval_period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The *X* minutes time period during which the **connection_fails_count** is evaluated",
		},
	}
}

func (me *ConnectionFailureDetection) MarshalHCL(properties hcl.Properties) error {
	if !me.Enabled {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"connection_fails_count": me.ConnectionFailsCount,
		"eval_period":            me.TimePeriodMinutes,
	})
}

func (me *ConnectionFailureDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("connection_fails_count"); ok {
		me.ConnectionFailsCount = opt.NewInt32(int32(value.(int)))
	}
	if value, ok := decoder.GetOk("eval_period"); ok {
		me.TimePeriodMinutes = opt.NewInt32(int32(value.(int)))
	}
	me.Enabled = true
	return nil
}
