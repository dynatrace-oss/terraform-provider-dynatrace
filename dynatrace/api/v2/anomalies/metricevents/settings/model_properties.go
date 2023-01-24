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

package metricevents

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ModelProperties struct {
	Type              ModelType      `json:"type"`                        // Metric-key-based query definitions only support static thresholds.
	Threshold         *float64       `json:"threshold,omitempty"`         // Raise an event if this value is violated
	AlertOnNoData     bool           `json:"alertOnNoData"`               // The ability to set an alert on missing data in a metric. When enabled, missing data samples will contribute as violating samples defined in advanced model properties. We recommend to not alert on missing data for sparse timeseries as this leads to alert spam.
	SignalFluctuation *float64       `json:"signalFluctuation,omitempty"` // Controls how many times the signal fluctuation is added to the baseline to produce the actual threshold for alerting
	Tolerance         *float64       `json:"tolerance,omitempty"`         //  Controls the width of the confidence band and larger values lead to a less sensitive model
	AlertCondition    AlertCondition `json:"alertCondition"`              // The alert condition of the model properties
	ViolatingSamples  int            `json:"violatingSamples"`            // The number of one-minute samples within the evaluation window that must violate to trigger an event.
	Samples           int            `json:"samples"`                     // The number of one-minute samples that form the sliding evaluation window.
	DealertingSamples int            `json:"dealertingSamples"`           // The number of one-minute samples within the evaluation window that must go back to normal to close the event.
}

func (me *ModelProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Metric-key-based query definitions only support static thresholds.",
			Required:    true,
		},
		"threshold": {
			Type:        schema.TypeFloat,
			Description: "Raise an event if this value is violated",
			Optional:    true,
		},
		"alert_on_no_data": {
			Type:        schema.TypeBool,
			Description: "The ability to set an alert on missing data in a metric. When enabled, missing data samples will contribute as violating samples defined in advanced model properties. We recommend to not alert on missing data for sparse timeseries as this leads to alert spam.",
			Required:    true,
		},
		"signal_fluctuation": {
			Type:        schema.TypeFloat,
			Description: "Controls how many times the signal fluctuation is added to the baseline to produce the actual threshold for alerting",
			Optional:    true,
		},
		"tolerance": {
			Type:        schema.TypeFloat,
			Description: "Controls the width of the confidence band and larger values lead to a less sensitive model",
			Optional:    true,
		},
		"alert_condition": {
			Type:        schema.TypeString,
			Description: "The alert condition of the model properties",
			Required:    true,
		},
		"violating_samples": {
			Type:        schema.TypeInt,
			Description: "The number of one-minute samples within the evaluation window that must violate to trigger an event.",
			Required:    true,
		},
		"samples": {
			Type:        schema.TypeInt,
			Description: "The number of one-minute samples that form the sliding evaluation window.",
			Required:    true,
		},
		"dealerting_samples": {
			Type:        schema.TypeInt,
			Description: "The number of one-minute samples within the evaluation window that must go back to normal to close the event.",
			Required:    true,
		},
	}
}

func (me *ModelProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":               me.Type,
		"threshold":          me.Threshold,
		"alert_on_no_data":   me.AlertOnNoData,
		"signal_fluctuation": me.SignalFluctuation,
		"tolerance":          me.Tolerance,
		"alert_condition":    me.AlertCondition,
		"violating_samples":  me.ViolatingSamples,
		"samples":            me.Samples,
		"dealerting_samples": me.DealertingSamples,
	})
}

func (me *ModelProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"type":               &me.Type,
		"threshold":          &me.Threshold,
		"alert_on_no_data":   &me.AlertOnNoData,
		"signal_fluctuation": &me.SignalFluctuation,
		"tolerance":          &me.Tolerance,
		"alert_condition":    &me.AlertCondition,
		"violating_samples":  &me.ViolatingSamples,
		"samples":            &me.Samples,
		"dealerting_samples": &me.DealertingSamples,
	})
	if me.Type == ModelTypes.Static && me.Threshold == nil {
		me.Threshold = opt.NewFloat64(0)
	}

	return err
}
