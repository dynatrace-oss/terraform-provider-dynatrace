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

package services

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No documentation available
type FailureRate struct {
	FixedDetection *FailureRateFixed `json:"fixedDetection,omitempty"` // . Alert if a given failure rate is exceeded during any 5-minute-period
	Enabled        bool              `json:"enabled"`                  // Detect increases in failure rate
	DetectionMode  *DetectionMode    `json:"detectionMode,omitempty"`  // Detection mode for increases in failure rate
	AutoDetection  *FailureRateAuto  `json:"autoDetection,omitempty"`  // . Alert if the percentage of failing service calls increases by **both** the absolute and relative thresholds:
}

func (me *FailureRate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fixed_detection": {
			Type:        schema.TypeList,
			Description: ". Alert if a given failure rate is exceeded during any 5-minute-period",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(FailureRateFixed).Schema()},
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Detect increases in failure rate",
			Required:    true,
		},
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Detection mode for increases in failure rate",
			Optional:    true,
		},
		"auto_detection": {
			Type:        schema.TypeList,
			Description: ". Alert if the percentage of failing service calls increases by **both** the absolute and relative thresholds:",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(FailureRateAuto).Schema()},
			Optional:    true,
		},
	}
}

func (me *FailureRate) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"fixed_detection": me.FixedDetection,
		"enabled":         me.Enabled,
		"detection_mode":  me.DetectionMode,
		"auto_detection":  me.AutoDetection,
	})
}

func (me *FailureRate) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"fixed_detection": &me.FixedDetection,
		"enabled":         &me.Enabled,
		"detection_mode":  &me.DetectionMode,
		"auto_detection":  &me.AutoDetection,
	})
}
