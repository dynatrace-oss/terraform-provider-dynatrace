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

package processgroups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AvailabilityMonitoringPG Configuration of the availability monitoring for the process group.
type AvailabilityMonitoring struct {
	Method           Method `json:"method"`                     // How to monitor the availability of the process group:  * `PROCESS_IMPACT`: Alert if any process of the group becomes unavailable.  * `MINIMUM_THRESHOLD`: Alert if the number of active processes in the group falls below the specified threshold.  * `OFF`: Availability monitoring is disabled.
	MinimumThreshold *int32 `json:"minimumThreshold,omitempty"` // Alert if the number of active processes in the group is lower than this value.
}

func (me *AvailabilityMonitoring) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"method": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "How to monitor the availability of the process group:  * `PROCESS_IMPACT`: Alert if any process of the group becomes unavailable.  * `MINIMUM_THRESHOLD`: Alert if the number of active processes in the group falls below the specified threshold.  * `OFF`: Availability monitoring is disabled.",
		},
		"minimum_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Alert if the number of active processes in the group is lower than this value.",
		},
	}
}

func (me *AvailabilityMonitoring) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]interface{}{
		"method":            me.Method,
		"minimum_threshold": me.MinimumThreshold,
	})
}

func (me *AvailabilityMonitoring) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]interface{}{
		"method":            &me.Method,
		"minimum_threshold": &me.MinimumThreshold,
	})
}

// Method How to monitor the availability of the process group:
// * `PROCESS_IMPACT`: Alert if any process of the group becomes unavailable.
// * `MINIMUM_THRESHOLD`: Alert if the number of active processes in the group falls below the specified threshold.
// * `OFF`: Availability monitoring is disabled.
type Method string

// Methods offers the known enum values
var Methods = struct {
	MinimumThreshold Method
	Off              Method
	ProcessImpact    Method
}{
	"MINIMUM_THRESHOLD",
	"OFF",
	"PROCESS_IMPACT",
}
