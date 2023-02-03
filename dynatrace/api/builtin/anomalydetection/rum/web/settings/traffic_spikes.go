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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TrafficSpikes struct {
	MinutesAbnormalState   float64 `json:"minutesAbnormalState"`   // Minutes an application has to stay in abnormal state before alert
	TrafficSpikePercentage float64 `json:"trafficSpikePercentage"` // Alert if the observed traffic is more than this percentage of the expected value
}

func (me *TrafficSpikes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"minutes_abnormal_state": {
			Type:        schema.TypeFloat,
			Description: "Minutes an application has to stay in abnormal state before alert",
			Required:    true,
		},
		"traffic_spike_percentage": {
			Type:        schema.TypeFloat,
			Description: "Alert if the observed traffic is more than this percentage of the expected value",
			Required:    true,
		},
	}
}

func (me *TrafficSpikes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"minutes_abnormal_state":   me.MinutesAbnormalState,
		"traffic_spike_percentage": me.TrafficSpikePercentage,
	})
}

func (me *TrafficSpikes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"minutes_abnormal_state":   &me.MinutesAbnormalState,
		"traffic_spike_percentage": &me.TrafficSpikePercentage,
	})
}
