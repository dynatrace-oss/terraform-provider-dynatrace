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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/load"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Autodetection Parameters of the response time degradation auto-detection. Required if the **detectionMode** is `DETECT_AUTOMATICALLY`. Not applicable otherwise.
// Violation of **any** criterion triggers an alert.
type Autodetection struct {
	LoadThreshold       load.Threshold             `json:"loadThreshold"`                              // Minimal service load to detect response time degradation.   Response time degradation of services with smaller load won't trigger alerts.
	Milliseconds        int32                      `json:"responseTimeDegradationMilliseconds"`        // Alert if the response time degrades by more than *X* milliseconds.
	Percent             int32                      `json:"responseTimeDegradationPercent"`             // Alert if the response time degrades by more than *X* %.
	SlowestMilliseconds int32                      `json:"slowestResponseTimeDegradationMilliseconds"` // Alert if the response time of the slowest 10% degrades by more than *X* milliseconds.
	SlowestPercent      int32                      `json:"slowestResponseTimeDegradationPercent"`      // Alert if the response time of the slowest 10% degrades by more than *X* %.
	Unknowns            map[string]json.RawMessage `json:"-"`
}

func (me *Autodetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load": {
			Type:        schema.TypeString,
			Description: "Minimal service load to detect response time degradation. Response time degradation of services with smaller load won't trigger alerts. Possible values are `FIFTEEN_REQUESTS_PER_MINUTE`, `FIVE_REQUESTS_PER_MINUTE`, `ONE_REQUEST_PER_MINUTE` and `TEN_REQUESTS_PER_MINUTE`",
			Required:    true,
		},
		"milliseconds": {
			Type:        schema.TypeInt,
			Description: "Alert if the response time degrades by more than *X* milliseconds",
			Required:    true,
		},
		"percent": {
			Type:        schema.TypeInt,
			Description: "Alert if the response time degrades by more than *X* %",
			Required:    true,
		},
		"slowest_milliseconds": {
			Type:        schema.TypeInt,
			Description: "Alert if the response time of the slowest 10% degrades by more than *X* milliseconds",
			Required:    true,
		},
		"slowest_percent": {
			Type:        schema.TypeInt,
			Description: "Alert if the response time of the slowest 10% degrades by more than *X* milliseconds",
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
		"load":                 me.LoadThreshold,
		"milliseconds":         me.Milliseconds,
		"percent":              me.Percent,
		"slowest_milliseconds": me.SlowestMilliseconds,
		"slowest_percent":      me.SlowestPercent,
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
		delete(me.Unknowns, "loadThreshold")
		delete(me.Unknowns, "responseTimeDegradationMilliseconds")
		delete(me.Unknowns, "responseTimeDegradationPercent")
		delete(me.Unknowns, "slowestResponseTimeDegradationMilliseconds")
		delete(me.Unknowns, "slowestResponseTimeDegradationPercent")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("load"); ok {
		me.LoadThreshold = load.Threshold(value.(string))
	}
	if value, ok := decoder.GetOk("milliseconds"); ok {
		me.Milliseconds = int32(value.(int))
	}
	if value, ok := decoder.GetOk("percent"); ok {
		me.Percent = int32(value.(int))
	}
	if value, ok := decoder.GetOk("slowest_milliseconds"); ok {
		me.SlowestMilliseconds = int32(value.(int))
	}
	if value, ok := decoder.GetOk("slowest_percent"); ok {
		me.SlowestPercent = int32(value.(int))
	}
	return nil
}

func (me *Autodetection) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"loadThreshold":                              me.LoadThreshold,
		"responseTimeDegradationMilliseconds":        me.Milliseconds,
		"responseTimeDegradationPercent":             me.Percent,
		"slowestResponseTimeDegradationMilliseconds": me.SlowestMilliseconds,
		"slowestResponseTimeDegradationPercent":      me.SlowestPercent,
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
		"loadThreshold":                              &me.LoadThreshold,
		"responseTimeDegradationMilliseconds":        &me.Milliseconds,
		"responseTimeDegradationPercent":             &me.Percent,
		"slowestResponseTimeDegradationMilliseconds": &me.SlowestMilliseconds,
		"slowestResponseTimeDegradationPercent":      &me.SlowestPercent,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
