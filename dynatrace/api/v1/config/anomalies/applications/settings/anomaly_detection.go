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

package applications

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/applications/settings/traffic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/detection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/failurerate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/responsetime"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AnomalyDetection The configuration of anomaly detection for applications.
type AnomalyDetection struct {
	ResponseTimeDegradation *responsetime.Detection `json:"responseTimeDegradation"` // Configuration of response time degradation detection.
	TrafficDrop             *traffic.DropDetection  `json:"trafficDrop"`             // The configuration of traffic drops detection.
	TrafficSpike            *traffic.SpikeDetection `json:"trafficSpike"`            // The configuration of traffic spikes detection.
	FailureRateIncrease     *failurerate.Detection  `json:"failureRateIncrease"`     // Configuration of failure rate increase detection.
}

func (me *AnomalyDetection) Name() string {
	return "application_anomalies"
}

func (me *AnomalyDetection) getFailureRateIncrease() *failurerate.Detection {
	if me.FailureRateIncrease == nil {
		return &failurerate.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	return me.FailureRateIncrease
}

func (me *AnomalyDetection) getResponseTimeDegradation() *responsetime.Detection {
	if me.ResponseTimeDegradation == nil {
		return &responsetime.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	return me.ResponseTimeDegradation
}

func (me *AnomalyDetection) GetTrafficSpike() *traffic.SpikeDetection {
	if me.TrafficSpike == nil {
		return &traffic.SpikeDetection{Enabled: false}
	}
	return me.TrafficSpike
}

func (me *AnomalyDetection) GetTrafficDrop() *traffic.DropDetection {
	if me.TrafficDrop == nil {
		return &traffic.DropDetection{Enabled: false}
	}
	return me.TrafficDrop
}

func (me *AnomalyDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"traffic": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for anomalies regarding traffic",
			Elem:        &schema.Resource{Schema: new(traffic.Detection).Schema()},
		},
		"response_time": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of response time degradation detection",
			Elem:        &schema.Resource{Schema: new(responsetime.Detection).Schema()},
		},
		"failure_rate": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of failure rate increase detection",
			Elem:        &schema.Resource{Schema: new(failurerate.Detection).Schema()},
		},
	}
}

func (me *AnomalyDetection) MarshalHCL(properties hcl.Properties) error {
	trafficDetection := &traffic.Detection{
		Drops:  me.TrafficDrop,
		Spikes: me.TrafficSpike,
	}
	properties.EncodeAll(map[string]any{
		"traffic":       trafficDetection,
		"response_time": me.ResponseTimeDegradation,
		"failure_rate":  me.FailureRateIncrease,
	})
	return nil
}

func (me *AnomalyDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("traffic.#"); ok {
		trafficDetection := new(traffic.Detection)
		if err := trafficDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "traffic", 0)); err != nil {
			return err
		}
		me.TrafficDrop = trafficDetection.Drops
		me.TrafficSpike = trafficDetection.Spikes
	}
	if _, ok := decoder.GetOk("response_time.#"); ok {
		me.ResponseTimeDegradation = new(responsetime.Detection)
		if err := me.ResponseTimeDegradation.UnmarshalHCL(hcl.NewDecoder(decoder, "response_time", 0)); err != nil {
			return err
		}
	} else {
		me.ResponseTimeDegradation = &responsetime.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	if _, ok := decoder.GetOk("failure_rate.#"); ok {
		me.FailureRateIncrease = new(failurerate.Detection)
		if err := me.FailureRateIncrease.UnmarshalHCL(hcl.NewDecoder(decoder, "failure_rate", 0)); err != nil {
			return err
		}
	} else {
		me.FailureRateIncrease = &failurerate.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	return nil
}

func (me *AnomalyDetection) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"trafficSpike":            me.GetTrafficSpike(),
		"responseTimeDegradation": me.getResponseTimeDegradation(),
		"failureRateIncrease":     me.getFailureRateIncrease(),
		"trafficDrop":             me.GetTrafficDrop(),
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *AnomalyDetection) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"trafficSpike":            &me.TrafficSpike,
		"responseTimeDegradation": &me.ResponseTimeDegradation,
		"failureRateIncrease":     &me.FailureRateIncrease,
		"trafficDrop":             &me.TrafficDrop,
	}); err != nil {
		return err
	}
	return nil
}
