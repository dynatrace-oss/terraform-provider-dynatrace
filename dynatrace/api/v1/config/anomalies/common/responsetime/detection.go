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

package responsetime

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/detection"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Detection Configuration of response time degradation detection.
type Detection struct {
	AutomaticDetection *Autodetection `json:"automaticDetection,omitempty"` // Parameters of the response time degradation auto-detection. Required if the **detectionMode** is `DETECT_AUTOMATICALLY`. Not applicable otherwise.  Violation of **any** criterion triggers an alert.
	DetectionMode      detection.Mode `json:"detectionMode"`                // How to detect response time degradation: automatically, or based on fixed thresholds, or do not detect.
	Thresholds         *Thresholds    `json:"thresholds,omitempty"`         // Fixed thresholds for response time degradation detection.   Required if **detectionMode** is `DETECT_USING_FIXED_THRESHOLDS`. Not applicable otherwise.
}

func (me *Detection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Parameters of the response time degradation auto-detection. Violation of **any** criterion triggers an alert",
			Elem:        &schema.Resource{Schema: new(Autodetection).Schema()},
		},
		"thresholds": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Fixed thresholds for response time degradation detection",
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
		me.AutomaticDetection = new(Autodetection)
		me.DetectionMode = detection.Modes.DetectAutomatically
		if err := me.AutomaticDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "auto", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("thresholds.#"); ok {
		me.Thresholds = new(Thresholds)
		me.DetectionMode = detection.Modes.DetectUsingFixedThresholds
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
		"automaticDetection": me.AutomaticDetection,
		"detectionMode":      me.DetectionMode,
		"thresholds":         me.Thresholds,
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
		"automaticDetection": &me.AutomaticDetection,
		"detectionMode":      &me.DetectionMode,
		"thresholds":         &me.Thresholds,
	}); err != nil {
		return err
	}
	return nil
}
