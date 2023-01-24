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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Autodetection Parameters of failure rate increase auto-detection. Required if **detectionMode** is `DETECT_AUTOMATICALLY`. Not applicable otherwise.
// The absolute and relative thresholds **both** must exceed to trigger an alert.
// Example: If the expected error rate is 1.5%, and you set an absolute increase of 1%, and a relative increase of 50%, the thresholds will be:
// Absolute: 1.5% + **1%** = 2.5%
// Relative: 1.5% + 1.5% * **50%** = 2.25%
type Autodetection struct {
	PercentAbsolute int32                      `json:"failingServiceCallPercentageIncreaseAbsolute"` // Absolute increase of failing service calls to trigger an alert, %.
	PercentRelative int32                      `json:"failingServiceCallPercentageIncreaseRelative"` // Relative increase of failing service calls to trigger an alert, %.
	Unknowns        map[string]json.RawMessage `json:"-"`
}

func (me *Autodetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"absolute": {
			Type:        schema.TypeInt,
			Description: "Absolute increase of failing service calls to trigger an alert, %",
			Required:    true,
		},
		"relative": {
			Type:        schema.TypeInt,
			Description: "Relative increase of failing service calls to trigger an alert, %",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Autodetection) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"absolute": me.PercentAbsolute,
		"relative": me.PercentRelative,
	})
}

func (me *Autodetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "failingServiceCallPercentageIncreaseAbsolute")
		delete(me.Unknowns, "failingServiceCallPercentageIncreaseRelative")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("relative"); ok {
		me.PercentRelative = int32(value.(int))
	}
	if value, ok := decoder.GetOk("absolute"); ok {
		me.PercentAbsolute = int32(value.(int))
	}
	return nil
}

func (me *Autodetection) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"failingServiceCallPercentageIncreaseAbsolute": me.PercentAbsolute,
		"failingServiceCallPercentageIncreaseRelative": me.PercentRelative,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Autodetection) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"failingServiceCallPercentageIncreaseAbsolute": &me.PercentAbsolute,
		"failingServiceCallPercentageIncreaseRelative": &me.PercentRelative,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
