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

package slo

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CustomDescription   *string              `json:"customDescription,omitempty"` // The description of the SLO
	Enabled             bool                 `json:"enabled"`                     // This setting is enabled (`true`) or disabled (`false`)
	ErrorBudgetBurnRate *ErrorBudgetBurnRate `json:"errorBudgetBurnRate"`         // ### Error budget burn rate
	EvaluationType      SloEvaluationType    `json:"evaluationType"`              // Possible Values: `AGGREGATE`
	EvaluationWindow    string               `json:"evaluationWindow"`            // Define the timeframe during which the SLO is to be evaluated. For the timeframe you can enter expressions like -1h (last hour), -1w (last week) or complex expressions like -2d to now (last two days), -1d/d to now/d (beginning of yesterday to beginning of today).
	Filter              string               `json:"filter"`                      // Set a filter parameter (entitySelector) on any GET call to evaluate this SLO against specific services only (for example, type(\"SERVICE\")).  For details, see the [Entity Selector documentation](https://dt-url.net/entityselector).
	MetricExpression    string               `json:"metricExpression"`            // For details, see the [Metrics page](/ui/metrics \"Metrics page\").
	MetricName          string               `json:"metricName"`                  // Metric name
	Name                string               `json:"name"`                        // SLO name
	TargetSuccess       float64              `json:"targetSuccess"`               // Set the target value of the SLO. A percentage below this value indicates a failure.
	TargetWarning       float64              `json:"targetWarning"`               // Set the warning value of the SLO. At the warning state the SLO is fulfilled. However, it is getting close to a failure.
	LegacyID            *string              `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_description": {
			Type:        schema.TypeString,
			Description: "The description of the SLO",
			Optional:    true, // nullable
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"error_budget_burn_rate": {
			Type:        schema.TypeList,
			Description: "### Error budget burn rate",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ErrorBudgetBurnRate).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"evaluation_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `AGGREGATE`",
			Required:    true,
		},
		"evaluation_window": {
			Type:        schema.TypeString,
			Description: "Define the timeframe during which the SLO is to be evaluated. For the timeframe you can enter expressions like -1h (last hour), -1w (last week) or complex expressions like -2d to now (last two days), -1d/d to now/d (beginning of yesterday to beginning of today).",
			Required:    true,
		},
		"filter": {
			Type:             schema.TypeString,
			Description:      "Set a filter parameter (entitySelector) on any GET call to evaluate this SLO against specific services only (for example, type(\"SERVICE\")).  For details, see the [Entity Selector documentation](https://dt-url.net/entityselector).",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressEOT,
		},
		"metric_expression": {
			Type:        schema.TypeString,
			Description: "For details, see the [Metrics page](/ui/metrics \"Metrics page\").",
			Required:    true,
		},
		"metric_name": {
			Type:             schema.TypeString,
			Description:      "Metric name",
			Optional:         true,
			Computed:         true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return d.Id() != "" },
		},
		"name": {
			Type:        schema.TypeString,
			Description: "SLO name",
			Required:    true,
		},
		"target_success": {
			Type:        schema.TypeFloat,
			Description: "Set the target value of the SLO. A percentage below this value indicates a failure.",
			Required:    true,
		},
		"target_warning": {
			Type:        schema.TypeFloat,
			Description: "Set the warning value of the SLO. At the warning state the SLO is fulfilled. However, it is getting close to a failure.",
			Required:    true,
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of this setting when referred to by the Config REST API V1",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"custom_description":     me.CustomDescription,
		"enabled":                me.Enabled,
		"error_budget_burn_rate": me.ErrorBudgetBurnRate,
		"evaluation_type":        me.EvaluationType,
		"evaluation_window":      me.EvaluationWindow,
		"filter":                 me.Filter,
		"metric_expression":      me.MetricExpression,
		"metric_name":            me.MetricName,
		"name":                   me.Name,
		"target_success":         me.TargetSuccess,
		"target_warning":         me.TargetWarning,
		"legacy_id":              me.LegacyID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_description":     &me.CustomDescription,
		"enabled":                &me.Enabled,
		"error_budget_burn_rate": &me.ErrorBudgetBurnRate,
		"evaluation_type":        &me.EvaluationType,
		"evaluation_window":      &me.EvaluationWindow,
		"filter":                 &me.Filter,
		"metric_expression":      &me.MetricExpression,
		"metric_name":            &me.MetricName,
		"name":                   &me.Name,
		"target_success":         &me.TargetSuccess,
		"target_warning":         &me.TargetWarning,
		"legacy_id":              &me.LegacyID,
	})
}
