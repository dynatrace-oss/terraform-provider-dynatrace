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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
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

type Objective struct {
	AutoAdaptiveThresholdEnabled *bool              `json:"autoAdaptiveThresholdEnabled,omitempty"` // Enable auto adaptive threshold
	ComparisonOperator           ComparisonOperator `json:"comparisonOperator"`                     // Possible Values: `GREATER_THAN_OR_EQUAL`, `LESS_THAN_OR_EQUAL`
	Description                  *string            `json:"description,omitempty"`
	DqlQuery                     *string            `json:"dqlQuery,omitempty"`     // DQL query
	Name                         string             `json:"name"`                   // Objective name
	ObjectiveType                ObjectiveType      `json:"objectiveType"`          // Possible Values: `DQL`, `REFERENCE_SLO`
	ReferenceSlo                 *string            `json:"referenceSlo,omitempty"` // Please enter the metric key of your desired SLO. SLO metric keys have to start with 'func:slo.'
	Target                       *float64           `json:"target,omitempty"`
	Warning                      *float64           `json:"warning,omitempty"`
}

func (me *Objective) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_adaptive_threshold_enabled": {
			Type:        schema.TypeBool,
			Description: "Enable auto adaptive threshold",
			Optional:    true, // precondition
		},
		"comparison_operator": {
			Type:        schema.TypeString,
			Description: "Possible Values: `GREATER_THAN_OR_EQUAL`, `LESS_THAN_OR_EQUAL`",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"dql_query": {
			Type:        schema.TypeString,
			Description: "DQL query",
			Optional:    true, // precondition
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Objective name",
			Required:    true,
		},
		"objective_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DQL`, `REFERENCE_SLO`",
			Required:    true,
		},
		"reference_slo": {
			Type:        schema.TypeString,
			Description: "Please enter the metric key of your desired SLO. SLO metric keys have to start with 'func:slo.'",
			Optional:    true, // precondition
		},
		"target": {
			Type:        schema.TypeFloat,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"warning": {
			Type:        schema.TypeFloat,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
	}
}

func (me *Objective) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_adaptive_threshold_enabled": me.AutoAdaptiveThresholdEnabled,
		"comparison_operator":             me.ComparisonOperator,
		"description":                     me.Description,
		"dql_query":                       me.DqlQuery,
		"name":                            me.Name,
		"objective_type":                  me.ObjectiveType,
		"reference_slo":                   me.ReferenceSlo,
		"target":                          me.Target,
		"warning":                         me.Warning,
	})
}

func (me *Objective) HandlePreconditions() error {
	if (me.AutoAdaptiveThresholdEnabled == nil) && (string(me.ObjectiveType) == "DQL") {
		me.AutoAdaptiveThresholdEnabled = opt.NewBool(false)
	}
	if (me.DqlQuery == nil) && (string(me.ObjectiveType) == "DQL") {
		return fmt.Errorf("'dql_query' must be specified if 'objective_type' is set to '%v'", me.ObjectiveType)
	}
	if (me.ReferenceSlo == nil) && (string(me.ObjectiveType) == "REFERENCE_SLO") {
		return fmt.Errorf("'reference_slo' must be specified if 'objective_type' is set to '%v'", me.ObjectiveType)
	}
	return nil
}

func (me *Objective) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_adaptive_threshold_enabled": &me.AutoAdaptiveThresholdEnabled,
		"comparison_operator":             &me.ComparisonOperator,
		"description":                     &me.Description,
		"dql_query":                       &me.DqlQuery,
		"name":                            &me.Name,
		"objective_type":                  &me.ObjectiveType,
		"reference_slo":                   &me.ReferenceSlo,
		"target":                          &me.Target,
		"warning":                         &me.Warning,
	})
}
