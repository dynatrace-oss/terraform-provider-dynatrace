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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Auto An auto-adaptive baseline strategy to detect anomalies within metrics that show a regular change over time, as the baseline is also updated automatically. An example is to detect an anomaly in the number of received network packets or within the number of user actions over time.
type Auto struct {
	BaseMonitoringStrategy
	AlertCondition             AlertCondition `json:"alertCondition"`                  // The condition for the **threshold** value check: above or below.
	AlertingOnMissingData      *bool          `json:"alertingOnMissingData,omitempty"` // If true, also one-minute samples without data are counted as violating samples.
	DealertingSamples          int32          `json:"dealertingSamples"`               // The number of one-minute samples within the evaluation window that must go back to normal to close the event.
	NumberOfSignalFluctuations float64        `json:"numberOfSignalFluctuations"`      // Defines the factor of how many signal fluctuations are valid. Values above the baseline plus the signal fluctuation times the number of tolerated signal fluctuations are alerted.
	Samples                    int32          `json:"samples"`                         // The number of one-minute samples that form the sliding evaluation window.
	ViolatingSamples           int32          `json:"violatingSamples"`                // The number of one-minute samples within the evaluation window that must violate the threshold to trigger an event.
}

func (me *Auto) GetType() Type {
	return Types.AutoAdaptiveBaseline
}

func (me *Auto) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alert_condition": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The condition for the **threshold** value check: `ABOVE` or `BELOW`",
		},
		"alerting_on_missing_data": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If true, also one-minute samples without data are counted as violating samples",
		},
		"dealerting_samples": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The number of one-minute samples within the evaluation window that must go back to normal to close the event",
		},
		"signal_fluctuations": {
			Type:        schema.TypeFloat,
			Required:    true,
			Description: "Defines the factor of how many signal fluctuations are valid. Values above the baseline plus the signal fluctuation times the number of tolerated signal fluctuations are alerted",
		},
		"samples": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The number of one-minute samples that form the sliding evaluation window",
		},
		"violating_samples": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The number of one-minute samples within the evaluation window that must violate the threshold to trigger an event",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Auto) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"alert_condition":          me.AlertCondition,
		"alerting_on_missing_data": me.AlertingOnMissingData,
		"dealerting_samples":       me.DealertingSamples,
		"signal_fluctuations":      me.NumberOfSignalFluctuations,
		"samples":                  me.Samples,
		"violating_samples":        me.ViolatingSamples,
	})
}

func (me *Auto) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "alertCondition")
		delete(me.Unknowns, "alertingOnMissingData")
		delete(me.Unknowns, "dealertingSamples")
		delete(me.Unknowns, "numberOfSignalFluctuations")
		delete(me.Unknowns, "samples")
		delete(me.Unknowns, "violatingSamples")
		delete(me.Unknowns, "type")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("alert_condition"); ok {
		me.AlertCondition = AlertCondition(value.(string))
	}
	if value, ok := decoder.GetOk("alerting_on_missing_data"); ok {
		me.AlertingOnMissingData = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("dealerting_samples"); ok {
		me.DealertingSamples = int32(value.(int))
	}
	if value, ok := decoder.GetOk("signal_fluctuations"); ok {
		me.NumberOfSignalFluctuations = value.(float64)
	}
	if value, ok := decoder.GetOk("samples"); ok {
		me.Samples = int32(value.(int))
	}
	if value, ok := decoder.GetOk("violating_samples"); ok {
		me.ViolatingSamples = int32(value.(int))
	}
	return nil
}

func (me *Auto) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"type":                       me.GetType(),
		"alertCondition":             me.AlertCondition,
		"alertingOnMissingData":      me.AlertingOnMissingData,
		"dealertingSamples":          me.DealertingSamples,
		"numberOfSignalFluctuations": me.NumberOfSignalFluctuations,
		"samples":                    me.Samples,
		"violatingSamples":           me.ViolatingSamples,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Auto) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"type":                       &me.Type,
		"alertCondition":             &me.AlertCondition,
		"alertingOnMissingData":      &me.AlertingOnMissingData,
		"dealertingSamples":          &me.DealertingSamples,
		"numberOfSignalFluctuations": &me.NumberOfSignalFluctuations,
		"samples":                    &me.Samples,
		"violatingSamples":           &me.ViolatingSamples,
	}); err != nil {
		return err
	}
	return nil
}
