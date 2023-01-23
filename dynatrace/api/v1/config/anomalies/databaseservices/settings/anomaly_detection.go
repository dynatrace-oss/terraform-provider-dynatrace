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

package databaseservices

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/detection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/failurerate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/load"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/responsetime"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AnomalyDetection The configuration of the anomaly detection for database services.
type AnomalyDetection struct {
	FailureRateIncrease            *failurerate.Detection      `json:"failureRateIncrease"`            // Configuration of failure rate increase detection.
	LoadDrop                       *load.DropDetection         `json:"loadDrop,omitempty"`             // The configuration of load drops detection.
	LoadSpike                      *load.SpikeDetection        `json:"loadSpike,omitempty"`            // The configuration of load spikes detection.
	ResponseTimeDegradation        *responsetime.Detection     `json:"responseTimeDegradation"`        // Configuration of response time degradation detection.
	DatabaseConnectionFailureCount *ConnectionFailureDetection `json:"databaseConnectionFailureCount"` // Parameters of the failed database connections detection.  The alert is triggered when failed connections number exceeds **connectionFailsCount** during any **timePeriodMinutes** minutes period.
}

func (me *AnomalyDetection) Name() string {
	return "database_anomalies"
}

func (me *AnomalyDetection) getDatabaseConnectionFailureCount() *ConnectionFailureDetection {
	if me.DatabaseConnectionFailureCount == nil {
		return &ConnectionFailureDetection{Enabled: false}
	}
	return me.DatabaseConnectionFailureCount
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

func (me *AnomalyDetection) getLoadSpike() *load.SpikeDetection {
	if me.LoadSpike == nil {
		return &load.SpikeDetection{Enabled: false}
	}
	return me.LoadSpike
}

func (me *AnomalyDetection) getLoadDrop() *load.DropDetection {
	if me.LoadDrop == nil {
		return &load.DropDetection{Enabled: false}
	}
	return me.LoadDrop
}

func (me *AnomalyDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for anomalies regarding load drops and spikes",
			Elem:        &schema.Resource{Schema: new(load.Detection).Schema()},
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
		"db_connect_failures": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Parameters of the failed database connections detection.  The alert is triggered when failed connections number exceeds **connectionFailsCount** during any **timePeriodMinutes** minutes period",
			Elem:        &schema.Resource{Schema: new(ConnectionFailureDetection).Schema()},
		},
	}
}

func (me *AnomalyDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"load": &load.Detection{
			Drops:  me.LoadDrop,
			Spikes: me.LoadSpike,
		},
		"response_time":       me.ResponseTimeDegradation,
		"failure_rate":        me.FailureRateIncrease,
		"db_connect_failures": me.DatabaseConnectionFailureCount,
	})
}

func (me *AnomalyDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("load.#"); ok {
		loadDetection := new(load.Detection)
		if err := loadDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "load", 0)); err != nil {
			return err
		}
		me.LoadDrop = loadDetection.Drops
		me.LoadSpike = loadDetection.Spikes
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
	if _, ok := decoder.GetOk("db_connect_failures.#"); ok {
		me.DatabaseConnectionFailureCount = new(ConnectionFailureDetection)
		if err := me.DatabaseConnectionFailureCount.UnmarshalHCL(hcl.NewDecoder(decoder, "db_connect_failures", 0)); err != nil {
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
		"loadSpike":                      me.getLoadSpike(),
		"responseTimeDegradation":        me.getResponseTimeDegradation(),
		"failureRateIncrease":            me.getFailureRateIncrease(),
		"loadDrop":                       me.getLoadDrop(),
		"databaseConnectionFailureCount": me.getDatabaseConnectionFailureCount(),
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
		"loadSpike":                      &me.LoadSpike,
		"responseTimeDegradation":        &me.ResponseTimeDegradation,
		"failureRateIncrease":            &me.FailureRateIncrease,
		"loadDrop":                       &me.LoadDrop,
		"databaseConnectionFailureCount": &me.DatabaseConnectionFailureCount,
	}); err != nil {
		return err
	}
	return nil
}
