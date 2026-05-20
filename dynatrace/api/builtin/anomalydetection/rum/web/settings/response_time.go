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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResponseTime struct {
	DetectionMode     *DetectionMode     `json:"detectionMode,omitempty"` // Detection strategy for key performance metric degradations. Possible values: `auto`, `fixed`
	Enabled           bool               `json:"enabled"`                 // This setting is enabled (`true`) or disabled (`false`)
	ResponseTimeAuto  *ResponseTimeAuto  `json:"responseTimeAuto,omitempty"`
	ResponseTimeFixed *ResponseTimeFixed `json:"responseTimeFixed,omitempty"`
}

func (me *ResponseTime) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Detection strategy for key performance metric degradations. Possible values: `auto`, `fixed`",
			Optional:    true, // precondition
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"response_time_auto": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ResponseTimeAuto).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"response_time_fixed": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ResponseTimeFixed).Schema()},
			MinItems:    1,
			MaxItems:    1,
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

func (me *ResponseTime) HandlePreconditions() error {
	if (me.DetectionMode != nil) && (!me.Enabled) {
		return fmt.Errorf("'detection_mode' must not be specified unless 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.DetectionMode == nil) && (me.Enabled) {
		return fmt.Errorf("'detection_mode' must be specified when 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.ResponseTimeAuto != nil) && (!me.Enabled || (me.DetectionMode == nil || string(*me.DetectionMode) != "auto")) {
		return fmt.Errorf("'response_time_auto' must not be specified unless ('enabled' is set to 'true' and 'detection_mode' is set to 'auto'); got 'enabled'='%v', 'detection_mode'='%v'", me.Enabled, me.DetectionMode)
	}
	if (me.ResponseTimeAuto == nil) && (me.Enabled && (me.DetectionMode != nil && string(*me.DetectionMode) == "auto")) {
		return fmt.Errorf("'response_time_auto' must be specified when ('enabled' is set to 'true' and 'detection_mode' is set to 'auto'); got 'enabled'='%v', 'detection_mode'='%v'", me.Enabled, me.DetectionMode)
	}
	if (me.ResponseTimeFixed != nil) && (!me.Enabled || (me.DetectionMode == nil || string(*me.DetectionMode) != "fixed")) {
		return fmt.Errorf("'response_time_fixed' must not be specified unless ('enabled' is set to 'true' and 'detection_mode' is set to 'fixed'); got 'enabled'='%v', 'detection_mode'='%v'", me.Enabled, me.DetectionMode)
	}
	if (me.ResponseTimeFixed == nil) && (me.Enabled && (me.DetectionMode != nil && string(*me.DetectionMode) == "fixed")) {
		return fmt.Errorf("'response_time_fixed' must be specified when ('enabled' is set to 'true' and 'detection_mode' is set to 'fixed'); got 'enabled'='%v', 'detection_mode'='%v'", me.Enabled, me.DetectionMode)
	}
	return nil
}

func (me *ResponseTime) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"detection_mode":      &me.DetectionMode,
		"enabled":             &me.Enabled,
		"response_time_auto":  &me.ResponseTimeAuto,
		"response_time_fixed": &me.ResponseTimeFixed,
	})
}
