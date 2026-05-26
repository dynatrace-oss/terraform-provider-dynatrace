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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ErrorRate struct {
	Enabled                bool            `json:"enabled"`                          // This setting is enabled (`true`) or disabled (`false`)
	ErrorRateAuto          *ErrorRateAuto  `json:"errorRateAuto,omitempty"`          // Alert if the percentage of failing user actions increases by **both** the absolute and relative thresholds:
	ErrorRateDetectionMode *DetectionMode  `json:"errorRateDetectionMode,omitempty"` // Detection strategy for increases in JavaScript errors. Possible values: `auto`, `fixed`
	ErrorRateFixed         *ErrorRateFixed `json:"errorRateFixed,omitempty"`
}

func (me *ErrorRate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"error_rate_auto": {
			Type:        schema.TypeList,
			Description: "Alert if the percentage of failing user actions increases by **both** the absolute and relative thresholds:",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ErrorRateAuto).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"error_rate_detection_mode": {
			Type:        schema.TypeString,
			Description: "Detection strategy for increases in JavaScript errors. Possible values: `auto`, `fixed`",
			Optional:    true, // precondition
		},
		"error_rate_fixed": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ErrorRateFixed).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *ErrorRate) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":                   me.Enabled,
		"error_rate_auto":           me.ErrorRateAuto,
		"error_rate_detection_mode": me.ErrorRateDetectionMode,
		"error_rate_fixed":          me.ErrorRateFixed,
	})
}

func (me *ErrorRate) HandlePreconditions() error {
	if (me.ErrorRateAuto != nil) && (!me.Enabled || (me.ErrorRateDetectionMode == nil || string(*me.ErrorRateDetectionMode) != "auto")) {
		return fmt.Errorf("'error_rate_auto' must not be specified unless ('enabled' is set to 'true' and 'error_rate_detection_mode' is set to 'auto'); got 'enabled'='%v', 'error_rate_detection_mode'='%v'", me.Enabled, opt.ValOrNil(me.ErrorRateDetectionMode))
	}
	if (me.ErrorRateAuto == nil) && (me.Enabled && (me.ErrorRateDetectionMode != nil && string(*me.ErrorRateDetectionMode) == "auto")) {
		return fmt.Errorf("'error_rate_auto' must be specified when ('enabled' is set to 'true' and 'error_rate_detection_mode' is set to 'auto'); got 'enabled'='%v', 'error_rate_detection_mode'='%v'", me.Enabled, opt.ValOrNil(me.ErrorRateDetectionMode))
	}
	if (me.ErrorRateDetectionMode != nil) && (!me.Enabled) {
		return fmt.Errorf("'error_rate_detection_mode' must not be specified unless 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.ErrorRateDetectionMode == nil) && (me.Enabled) {
		return fmt.Errorf("'error_rate_detection_mode' must be specified when 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.ErrorRateFixed != nil) && (!me.Enabled || (me.ErrorRateDetectionMode == nil || string(*me.ErrorRateDetectionMode) != "fixed")) {
		return fmt.Errorf("'error_rate_fixed' must not be specified unless ('enabled' is set to 'true' and 'error_rate_detection_mode' is set to 'fixed'); got 'enabled'='%v', 'error_rate_detection_mode'='%v'", me.Enabled, opt.ValOrNil(me.ErrorRateDetectionMode))
	}
	if (me.ErrorRateFixed == nil) && (me.Enabled && (me.ErrorRateDetectionMode != nil && string(*me.ErrorRateDetectionMode) == "fixed")) {
		return fmt.Errorf("'error_rate_fixed' must be specified when ('enabled' is set to 'true' and 'error_rate_detection_mode' is set to 'fixed'); got 'enabled'='%v', 'error_rate_detection_mode'='%v'", me.Enabled, opt.ValOrNil(me.ErrorRateDetectionMode))
	}
	return nil
}

func (me *ErrorRate) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":                   &me.Enabled,
		"error_rate_auto":           &me.ErrorRateAuto,
		"error_rate_detection_mode": &me.ErrorRateDetectionMode,
		"error_rate_fixed":          &me.ErrorRateFixed,
	})
}
