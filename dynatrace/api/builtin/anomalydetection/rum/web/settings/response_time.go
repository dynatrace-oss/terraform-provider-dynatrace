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

type ResponseTime struct {
	DetectionMode     *DetectionMode     `json:"detectionMode,omitempty"` // Possible Values: `Auto`, `Fixed`
	Enabled           bool               `json:"enabled"`                 // This setting is enabled (`true`) or disabled (`false`)
	ResponseTimeAuto  *ResponseTimeAuto  `json:"responseTimeAuto,omitempty"`
	ResponseTimeFixed *ResponseTimeFixed `json:"responseTimeFixed,omitempty"`
}

func (me *ResponseTime) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Auto`, `Fixed`",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"response_time_auto": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(ResponseTimeAuto).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"response_time_fixed": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(ResponseTimeFixed).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *ResponseTime) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"detection_mode":      me.DetectionMode,
		"enabled":             me.Enabled,
		"response_time_auto":  me.ResponseTimeAuto,
		"response_time_fixed": me.ResponseTimeFixed,
	})
}

func (me *ResponseTime) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"detection_mode":      &me.DetectionMode,
		"enabled":             &me.Enabled,
		"response_time_auto":  &me.ResponseTimeAuto,
		"response_time_fixed": &me.ResponseTimeFixed,
	})
}
