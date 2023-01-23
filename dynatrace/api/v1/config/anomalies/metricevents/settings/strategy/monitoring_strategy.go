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

package strategy

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MonitoringStrategy A monitoring strategy for a metric event config.
// This is the base version of the monitoring strategy, depending on the type,
// the actual JSON may contain additional fields.
type MonitoringStrategy interface {
	GetType() Type
}

// BaseMetricEventMonitoringStrategy A monitoring strategy for a metric event config.
// This is the base version of the monitoring strategy, depending on the type,
// the actual JSON may contain additional fields.
type BaseMonitoringStrategy struct {
	Type     Type                       `json:"type"` // Defines the actual set of fields depending on the value. See one of the following objects:  * `STATIC_THRESHOLD` -> MetricEventStaticThresholdMonitoringStrategy  * `AUTO_ADAPTIVE_BASELINE` -> MetricEventAutoAdaptiveBaselineMonitoringStrategy
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *BaseMonitoringStrategy) GetType() Type {
	return me.Type
}

func (me *BaseMonitoringStrategy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the actual set of fields depending on the value",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *BaseMonitoringStrategy) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.Encode("type", me.Type)
}

func (me *BaseMonitoringStrategy) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "type")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = Type(value.(string))
	}
	return nil
}

func (me *BaseMonitoringStrategy) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"type": me.Type,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *BaseMonitoringStrategy) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"type": &me.Type,
	}); err != nil {
		return err
	}
	return nil
}

// Type Defines the actual set of fields depending on the value. See one of the following objects:
// * `STATIC_THRESHOLD` -> MetricEventStaticThresholdMonitoringStrategy
// * `AUTO_ADAPTIVE_BASELINE` -> MetricEventAutoAdaptiveBaselineMonitoringStrategy
type Type string

// Types offers the known enum values
var Types = struct {
	AutoAdaptiveBaseline Type
	StaticThreshold      Type
}{
	"AUTO_ADAPTIVE_BASELINE",
	"STATIC_THRESHOLD",
}

// AlertCondition The condition for the **threshold** value check: above or below.
type AlertCondition string

// AlertConditions offers the known enum values
var AlertConditions = struct {
	Above AlertCondition
	Below AlertCondition
}{
	"ABOVE",
	"BELOW",
}
