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
type ResponseTime struct {
	AutoDetection  *ResponseTimeAuto  `json:"autoDetection,omitempty"`  // No documentation available
	FixedDetection *ResponseTimeFixed `json:"fixedDetection,omitempty"` // No documentation available
	Enabled        bool               `json:"enabled"`                  // Detect response time degradations
	DetectionMode  *DetectionMode     `json:"detectionMode,omitempty"`  // Detection mode for response time degradations
}

func (me *ResponseTime) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_detection": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ResponseTimeAuto).Schema()},
			Optional:    true,
		},
		"fixed_detection": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ResponseTimeFixed).Schema()},
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Detect response time degradations",
			Required:    true,
		},
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Detection mode for response time degradations",
			Optional:    true,
		},
	}
}

func (me *ResponseTime) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_detection":  me.AutoDetection,
		"fixed_detection": me.FixedDetection,
		"enabled":         me.Enabled,
		"detection_mode":  me.DetectionMode,
	})
}

func (me *ResponseTime) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_detection":  &me.AutoDetection,
		"fixed_detection": &me.FixedDetection,
		"enabled":         &me.Enabled,
		"detection_mode":  &me.DetectionMode,
	})
}
