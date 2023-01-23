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

package failurerate

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/detection"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Detection Configuration of failure rate increase detection.
type Detection struct {
	Thresholds         *Thresholds    `json:"thresholds,omitempty"`         // Fixed thresholds for failure rate increase detection.   Required if **detectionMode** is `DETECT_USING_FIXED_THRESHOLDS`. Not applicable otherwise.
	AutomaticDetection *Autodetection `json:"automaticDetection,omitempty"` // Parameters of failure rate increase auto-detection. Required if **detectionMode** is `DETECT_AUTOMATICALLY`. Not applicable otherwise.  The absolute and relative thresholds **both** must exceed to trigger an alert.  Example: If the expected error rate is 1.5%, and you set an absolute increase of 1%, and a relative increase of 50%, the thresholds will be:  Absolute: 1.5% + **1%** = 2.5%  Relative: 1.5% + 1.5% * **50%** = 2.25%
	DetectionMode      detection.Mode `json:"detectionMode"`                // How to detect failure rate increase: automatically, or based on fixed thresholds, or do not detect.
}

func (me *Detection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Parameters of failure rate increase auto-detection. Example: If the expected error rate is 1.5%, and you set an absolute increase of 1%, and a relative increase of 50%, the thresholds will be:  Absolute: 1.5% + **1%** = 2.5%  Relative: 1.5% + 1.5% * **50%** = 2.25%",
			Elem:        &schema.Resource{Schema: new(Autodetection).Schema()},
		},
		"thresholds": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Fixed thresholds for failure rate increase detection",
			Elem:        &schema.Resource{Schema: new(Thresholds).Schema()},
		},
	}
}

func (me *Detection) MarshalHCL(properties hcl.Properties) error {
	if me.DetectionMode == detection.Modes.DontDetect {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"auto":       me.AutomaticDetection,
		"thresholds": me.Thresholds,
	})
}

func (me *Detection) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("auto.#"); ok {
		me.DetectionMode = detection.Modes.DetectAutomatically
		me.AutomaticDetection = new(Autodetection)
		if err := me.AutomaticDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "auto", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("thresholds.#"); ok {
		me.DetectionMode = detection.Modes.DetectUsingFixedThresholds
		me.Thresholds = new(Thresholds)
		if err := me.Thresholds.UnmarshalHCL(hcl.NewDecoder(decoder, "thresholds", 0)); err != nil {
			return err
		}
	} else {
		me.DetectionMode = detection.Modes.DontDetect
	}
	return nil
}

func (me *Detection) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"thresholds":         me.Thresholds,
		"automaticDetection": me.AutomaticDetection,
		"detectionMode":      me.DetectionMode,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Detection) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"thresholds":         &me.Thresholds,
		"automaticDetection": &me.AutomaticDetection,
		"detectionMode":      &me.DetectionMode,
	}); err != nil {
		return err
	}
	return nil
}
