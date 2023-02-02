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

type HighNetworkDetectionThresholds struct {
	ErrorsPercentage int              `json:"errorsPercentage"` // Alert if sent/received traffic utilization is higher than this threshold for the defined amount of samples
	EventThresholds  *EventThresholds `json:"eventThresholds"`
}

func (me *HighNetworkDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"errors_percentage": {
			Type:        schema.TypeInt,
			Description: "Alert if sent/received traffic utilization is higher than this threshold for the defined amount of samples",
			Required:    true,
		},
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *HighNetworkDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"errors_percentage": me.ErrorsPercentage,
		"event_thresholds":  me.EventThresholds,
	})
}

func (me *HighNetworkDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"errors_percentage": &me.ErrorsPercentage,
		"event_thresholds":  &me.EventThresholds,
	})
}
