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

package infrastructurehosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NetworkHighRetransmissionDetectionThresholds struct {
	EventThresholds                     *EventThresholds `json:"eventThresholds"`
	RetransmissionRatePercentage        int              `json:"retransmissionRatePercentage"`        // Retransmission rate threshold
	RetransmittedPacketsNumberPerMinute int              `json:"retransmittedPacketsNumberPerMinute"` // Number of retransmitted packets threshold
}

func (me *NetworkHighRetransmissionDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"retransmission_rate_percentage": {
			Type:        schema.TypeInt,
			Description: "Retransmission rate threshold",
			Required:    true,
		},
		"retransmitted_packets_number_per_minute": {
			Type:        schema.TypeInt,
			Description: "Number of retransmitted packets threshold",
			Required:    true,
		},
	}
}

func (me *NetworkHighRetransmissionDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_thresholds":                        me.EventThresholds,
		"retransmission_rate_percentage":          me.RetransmissionRatePercentage,
		"retransmitted_packets_number_per_minute": me.RetransmittedPacketsNumberPerMinute,
	})
}

func (me *NetworkHighRetransmissionDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_thresholds":                        &me.EventThresholds,
		"retransmission_rate_percentage":          &me.RetransmissionRatePercentage,
		"retransmitted_packets_number_per_minute": &me.RetransmittedPacketsNumberPerMinute,
	})
}
