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

package sitereliabilityguardian

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Objectives []*Objective

func (me *Objectives) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"objective": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Objective).Schema()},
		},
	}
}

func (me Objectives) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("objective", me)
}

func (me *Objectives) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("objective", me)
}

// Objective. A single validation criterion evaluated each guardian run. Result is pass, warning, fail, or info depending on thresholds and the comparison operator.
type Objective struct {
	AutoAdaptiveThresholdEnabled *bool              `json:"autoAdaptiveThresholdEnabled,omitempty"` // Dynamically computes thresholds from 30 days of history.
	ComparisonOperator           ComparisonOperator `json:"comparisonOperator"`                     // Pass/fail direction: use ≥ when higher values are better, ≤ when lower values are better. Possible values: `GREATER_THAN_OR_EQUAL`, `LESS_THAN_OR_EQUAL`
	Description                  *string            `json:"description,omitempty"`                  // Optional short explanation of what this objective measures.
	DisplayUnit                  *DisplayUnit       `json:"displayUnit,omitempty"`                  // Optional unit conversion and decimal formatting applied when displaying the DQL result in the UI.
	DqlQuery                     *string            `json:"dqlQuery,omitempty"`                     // DQL query to execute. The first numeric result becomes the objective value. Supports $variable interpolation.
	Links                        ObjectiveLinks     `json:"links,omitempty"`                        // Fields for adding relevant links to this objective.
	Name                         string             `json:"name"`                                   // Unique name within this guardian. Included in every emitted validation event as the objective identifier.
	ObjectiveType                ObjectiveType      `json:"objectiveType"`                          // How the objective value is computed: via a DQL query or an existing SLO metric. Possible values: `DQL`, `REFERENCE_SLO`
	ReferenceSlo                 *string            `json:"referenceSlo,omitempty"`                 // Please enter the metric key of your desired SLO. SLO metric keys have to start with 'func:slo.'
	Segments                     Segments           `json:"segments,omitempty"`                     // Optional Grail segments to scope the DQL query to specific data.
	Target                       *float64           `json:"target,omitempty"`                       // Hard pass/fail threshold. Missing this value yields FAIL. If unset with no warning, status is always INFO.
	Warning                      *float64           `json:"warning,omitempty"`                      // Soft threshold. Results between warning and target yield WARNING. When set alone, yields PASS or WARNING.
}

func (me *Objective) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_adaptive_threshold_enabled": {
			Type:        schema.TypeBool,
			Description: "Dynamically computes thresholds from 30 days of history.",
			Optional:    true, // nullable & precondition
		},
		"comparison_operator": {
			Type:        schema.TypeString,
			Description: "Pass/fail direction: use ≥ when higher values are better, ≤ when lower values are better. Possible values: `GREATER_THAN_OR_EQUAL`, `LESS_THAN_OR_EQUAL`",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Optional short explanation of what this objective measures.",
			Optional:    true, // nullable
		},
		"display_unit": {
			Type:        schema.TypeList,
			Description: "Optional unit conversion and decimal formatting applied when displaying the DQL result in the UI.",
			Optional:    true, // nullable & precondition
			Elem:        &schema.Resource{Schema: new(DisplayUnit).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"dql_query": {
			Type:        schema.TypeString,
			Description: "DQL query to execute. The first numeric result becomes the objective value. Supports $variable interpolation.",
			Optional:    true, // precondition
		},
		"links": {
			Type:        schema.TypeList,
			Description: "Fields for adding relevant links to this objective.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(ObjectiveLinks).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Unique name within this guardian. Included in every emitted validation event as the objective identifier.",
			Required:    true,
		},
		"objective_type": {
			Type:        schema.TypeString,
			Description: "How the objective value is computed: via a DQL query or an existing SLO metric. Possible values: `DQL`, `REFERENCE_SLO`",
			Required:    true,
		},
		"reference_slo": {
			Type:        schema.TypeString,
			Description: "Please enter the metric key of your desired SLO. SLO metric keys have to start with 'func:slo.'",
			Optional:    true, // precondition
		},
		"segments": {
			Type:        schema.TypeList,
			Description: "Optional Grail segments to scope the DQL query to specific data.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Segments).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"target": {
			Type:        schema.TypeFloat,
			Description: "Hard pass/fail threshold. Missing this value yields FAIL. If unset with no warning, status is always INFO.",
			Optional:    true, // nullable
		},
		"warning": {
			Type:        schema.TypeFloat,
			Description: "Soft threshold. Results between warning and target yield WARNING. When set alone, yields PASS or WARNING.",
			Optional:    true, // nullable
		},
	}
}

func (me *Objective) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_adaptive_threshold_enabled": me.AutoAdaptiveThresholdEnabled,
		"comparison_operator":             me.ComparisonOperator,
		"description":                     me.Description,
		"display_unit":                    me.DisplayUnit,
		"dql_query":                       me.DqlQuery,
		"links":                           me.Links,
		"name":                            me.Name,
		"objective_type":                  me.ObjectiveType,
		"reference_slo":                   me.ReferenceSlo,
		"segments":                        me.Segments,
		"target":                          me.Target,
		"warning":                         me.Warning,
	})
}

func (me *Objective) HandlePreconditions() error {
	if (me.AutoAdaptiveThresholdEnabled == nil) && (string(me.ObjectiveType) == "DQL") {
		me.AutoAdaptiveThresholdEnabled = new(false)
	}
	if (me.AutoAdaptiveThresholdEnabled != nil) && (string(me.ObjectiveType) != "DQL") {
		return fmt.Errorf("'auto_adaptive_threshold_enabled' must not be specified unless 'objective_type' is set to 'DQL'; got 'objective_type'='%v'", me.ObjectiveType)
	}
	if (me.DisplayUnit != nil) && (string(me.ObjectiveType) != "DQL") {
		return fmt.Errorf("'display_unit' must not be specified unless 'objective_type' is set to 'DQL'; got 'objective_type'='%v'", me.ObjectiveType)
	}
	if (me.DqlQuery != nil) && (string(me.ObjectiveType) != "DQL") {
		return fmt.Errorf("'dql_query' must not be specified unless 'objective_type' is set to 'DQL'; got 'objective_type'='%v'", me.ObjectiveType)
	}
	if (me.DqlQuery == nil) && (string(me.ObjectiveType) == "DQL") {
		return fmt.Errorf("'dql_query' must be specified when 'objective_type' is set to 'DQL'; got 'objective_type'='%v'", me.ObjectiveType)
	}
	if (me.ReferenceSlo != nil) && (string(me.ObjectiveType) != "REFERENCE_SLO") {
		return fmt.Errorf("'reference_slo' must not be specified unless 'objective_type' is set to 'REFERENCE_SLO'; got 'objective_type'='%v'", me.ObjectiveType)
	}
	if (me.ReferenceSlo == nil) && (string(me.ObjectiveType) == "REFERENCE_SLO") {
		return fmt.Errorf("'reference_slo' must be specified when 'objective_type' is set to 'REFERENCE_SLO'; got 'objective_type'='%v'", me.ObjectiveType)
	}
	return nil
}

func (me *Objective) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_adaptive_threshold_enabled": &me.AutoAdaptiveThresholdEnabled,
		"comparison_operator":             &me.ComparisonOperator,
		"description":                     &me.Description,
		"display_unit":                    &me.DisplayUnit,
		"dql_query":                       &me.DqlQuery,
		"links":                           &me.Links,
		"name":                            &me.Name,
		"objective_type":                  &me.ObjectiveType,
		"reference_slo":                   &me.ReferenceSlo,
		"segments":                        &me.Segments,
		"target":                          &me.Target,
		"warning":                         &me.Warning,
	})
}
