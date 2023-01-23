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
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SLO struct {
	Name              string  `json:"name"`                        // The name of the SLO
	Enabled           bool    `json:"enabled,omitempty"`           // The SLO is enabled (`true`) or disabled (`false`)
	Description       *string `json:"customDescription,omitempty"` // The custom description of the SLO (optional)
	UseRateMetric     bool    `json:"useRateMetric"`               // The type of the metric to use for SLO calculation: \n\n* `true`: An existing percentage-based metric. \n* `false`: A ratio of two metrics. \n\nFor a list of available metrics, see [Built-in metric page](https://dt-url.net/be03kow) or try the [GET metrics](https://dt-url.net/8e43kxf) API call
	MetricRate        *string `json:"metricRate,omitempty"`        // The percentage-based metric for the calculation of the SLO. \n\nRequired when the **useRateMetric** is set to `true`
	MetricExpression  *string `json:"metricExpression,omitempty"`  // The percentage-based metric expression for the calculation of the SLO
	MetricNumerator   *string `json:"metricNumerator,omitempty"`   // The metric for the count of successes (the numerator in rate calculation). \n\nRequired when the **useRateMetric** is set to `false`
	MetricDenominator *string `json:"metricDenominator,omitempty"` // The total count metric (the denominator in rate calculation). \n\nRequired when the **useRateMetric** is set to `false`
	EvaluationType    string  `json:"evaluationType"`              // The evaluation type of the SLO. Currently only `AGGREGATE` is supported
	Filter            *string `json:"filter,omitempty"`            // The entity filter for the SLO evaluation. Use the [syntax of entity selector](https://dt-url.net/entityselector)
	Target            float64 `json:"target"`                      // The target value of the SLO
	Warning           float64 `json:"warning"`                     // The warning value of the SLO. \n\n At warning state the SLO is still fulfilled but is getting close to failure
	Timeframe         string  `json:"timeframe"`                   // The timeframe for the SLO evaluation. Use the syntax of the global timeframe selector
}

func (me *SLO) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the rule",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The custom description of the SLO (optional)",
		},
		"metric_expression": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The percentage-based metric expression for the calculation of the SLO",
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
		"disabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The SLO is enabled (`false`) or disabled (`true`)",
		},
		"rate": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"numerator", "denominator"},
			Description:   "The percentage-based metric for the calculation of the SLO",
		},
		"numerator": {
			Type:          schema.TypeString,
			Optional:      true,
			RequiredWith:  []string{"denominator"},
			ConflictsWith: []string{"rate"},
			Deprecated:    "`numerator` and `denominator` have been replaced by `metric_expression`",
			Description:   "The metric for the count of successes (the numerator in rate calculation)",
		},
		"denominator": {
			Type:          schema.TypeString,
			Optional:      true,
			RequiredWith:  []string{"numerator"},
			ConflictsWith: []string{"rate"},
			Description:   "The total count metric (the denominator in rate calculation)",
		},
		"evaluation": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The evaluation type of the SLO. Currently only `AGGREGATE` is supported",
		},
		"filter": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The entity filter for the SLO evaluation. Use the [syntax of entity selector](https://dt-url.net/entityselector)",
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
		"target": {
			Type:        schema.TypeFloat,
			Required:    true,
			Description: "The target value of the SLO",
		},
		"warning": {
			Type:        schema.TypeFloat,
			Required:    true,
			Description: "The warning value of the SLO. At warning state the SLO is still fulfilled but is getting close to failure",
		},
		"timeframe": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The timeframe for the SLO evaluation. Use the syntax of the global timeframe selector",
		},
	}
}

func empty2Nil(s *string) *string {
	if s == nil {
		return nil
	}
	if len(*s) == 0 {
		return nil
	}
	return s
}

func (me *SLO) MarshalHCL(properties hcl.Properties) error {

	err := properties.EncodeAll(map[string]any{
		"name":              me.Name,
		"description":       me.Description,
		"disabled":          !me.Enabled,
		"rate":              empty2Nil(me.MetricRate),
		"metric_expression": empty2Nil(me.MetricExpression),
		"numerator":         empty2Nil(me.MetricNumerator),
		"denominator":       empty2Nil(me.MetricDenominator),
		"evaluation":        me.EvaluationType,
		"filter":            me.Filter,
		"target":            me.Target,
		"warning":           me.Warning,
		"timeframe":         me.Timeframe,
	})
	if err != nil {
		return err
	}
	if me.Enabled {
		delete(properties, "disabled")
	}
	if me.MetricRate != nil {
		var mr = *me.MetricRate
		if strings.HasSuffix(mr, ":splitBy():splitBy()") {
			mr = mr[0 : len(mr)-len(":splitBy()")]
			me.MetricRate = &mr
			if err := properties.Encode("rate", mr); err != nil {
				return err
			}
		}
	}
	if me.MetricNumerator != nil {
		var mr = *me.MetricNumerator
		if strings.HasSuffix(mr, ":splitBy():splitBy()") {
			mr = mr[0 : len(mr)-len(":splitBy()")]
			me.MetricNumerator = &mr
			if err := properties.Encode("numerator", mr); err != nil {
				return err
			}
		}
	}
	if me.MetricDenominator != nil {
		var mr = *me.MetricDenominator
		if strings.HasSuffix(mr, ":splitBy():splitBy()") {
			mr = mr[0 : len(mr)-len(":splitBy()")]
			me.MetricDenominator = &mr
			if err := properties.Encode("denominator", mr); err != nil {
				return err
			}
		}
	}
	return nil
}

func nonNil(s *string) *string {
	if s == nil {
		return opt.NewString("")
	}
	return s
}

func (me *SLO) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"name":              &me.Name,
		"description":       &me.Description,
		"disabled":          &me.Enabled,
		"rate":              &me.MetricRate,
		"numerator":         &me.MetricNumerator,
		"metric_expression": &me.MetricExpression,
		"denominator":       &me.MetricDenominator,
		"evaluation":        &me.EvaluationType,
		"filter":            &me.Filter,
		"target":            &me.Target,
		"warning":           &me.Warning,
		"timeframe":         &me.Timeframe,
	})
	me.Enabled = !me.Enabled
	me.MetricNumerator = nonNil(me.MetricNumerator)
	me.MetricDenominator = nonNil(me.MetricDenominator)
	me.MetricRate = nonNil(me.MetricRate)
	// me.MetricExpression = nonNil(me.MetricExpression)
	me.UseRateMetric = (me.MetricRate != nil) && len(*me.MetricRate) > 0
	return err
}
