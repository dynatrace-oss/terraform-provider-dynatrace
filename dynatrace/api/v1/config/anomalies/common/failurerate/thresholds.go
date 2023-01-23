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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Thresholds Fixed thresholds for failure rate increase detection.
//
//	Required if **detectionMode** is `DETECT_USING_FIXED_THRESHOLDS`. Not applicable otherwise.
type Thresholds struct {
	Sensitivity common.Sensitivity         `json:"sensitivity"` // Sensitivity of the threshold.  With `low` sensitivity, high statistical confidence is used. Brief violations (for example, due to a surge in load) won't trigger alerts.  With `high` sensitivity, no statistical confidence is used. Each violation triggers alert.
	Threshold   int32                      `json:"threshold"`   // Failure rate during any 5-minute period to trigger an alert, %.
	Unknowns    map[string]json.RawMessage `json:"-"`
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sensitivity": {
			Type:        schema.TypeString,
			Description: "Sensitivity of the threshold.  With `low` sensitivity, high statistical confidence is used. Brief violations (for example, due to a surge in load) won't trigger alerts.  With `high` sensitivity, no statistical confidence is used. Each violation triggers alert",
			Required:    true,
		},
		"threshold": {
			Type:        schema.TypeInt,
			Description: "Failure rate during any 5-minute period to trigger an alert, %",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Thresholds) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"sensitivity": me.Sensitivity,
		"threshold":   me.Threshold,
	})
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "sensitivity")
		delete(me.Unknowns, "threshold")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("sensitivity"); ok {
		me.Sensitivity = common.Sensitivity(value.(string))
	}
	if value, ok := decoder.GetOk("threshold"); ok {
		me.Threshold = int32(value.(int))
	}
	return nil
}

func (me *Thresholds) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"sensitivity": me.Sensitivity,
		"threshold":   me.Threshold,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Thresholds) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"sensitivity": &me.Sensitivity,
		"threshold":   &me.Threshold,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
