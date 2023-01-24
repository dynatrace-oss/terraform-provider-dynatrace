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

package service

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DimensionDefinition Parameters of a definition of a calculated service metric.
type DimensionDefinition struct {
	TopXDirection   TopXDirection              `json:"topXDirection"`          // How to calculate the **topX** values.
	Dimension       string                     `json:"dimension"`              // The dimension value pattern.   You can define custom placeholders in the **placeholders** field and use them here.
	Name            string                     `json:"name"`                   // The name of the dimension.
	Placeholders    Placeholders               `json:"placeholders,omitempty"` // The list of custom placeholders to be used in a dimension value pattern.
	TopX            int32                      `json:"topX"`                   // The number of top values to be calculated.
	TopXAggregation TopXAggregation            `json:"topXAggregation"`        // The aggregation of the dimension.
	Unknowns        map[string]json.RawMessage `json:"-"`
}

func (me *DimensionDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the dimension",
		},
		"dimension": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The dimension value pattern. You can define custom placeholders in the `placeholders` field and use them here",
		},
		"top_x_direction": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "How to calculate the **topX** values. Possible values are `ASCENDING` and `DESCENDING`",
		},
		"top_x": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The number of top values to be calculated",
		},
		"top_x_aggregation": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The aggregation of the dimension. Possible values are `AVERAGE`, `COUNT`, `MAX`, `MIN`, `OF_INTEREST_RATIO`, `OTHER_RATIO`, `SINGLE_VALUE` and `SUM`",
		},
		"placeholders": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "The list of custom placeholders to be used in a dimension value pattern",
			Elem:        &schema.Resource{Schema: new(Placeholders).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *DimensionDefinition) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"name":              me.Name,
		"dimension":         me.Dimension,
		"top_x_direction":   me.TopXDirection,
		"top_x":             me.TopX,
		"top_x_aggregation": me.TopXAggregation,
		"placeholders":      me.Placeholders,
		"unknowns":          me.Unknowns,
	})
}

func (me *DimensionDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":              &me.Name,
		"dimension":         &me.Dimension,
		"top_x_direction":   &me.TopXDirection,
		"top_x":             &me.TopX,
		"top_x_aggregation": &me.TopXAggregation,
		"placeholders":      &me.Placeholders,
		"unknowns":          &me.Unknowns,
	})
}

func (me *DimensionDefinition) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"name":            me.Name,
		"dimension":       me.Dimension,
		"placeholders":    me.Placeholders,
		"topXDirection":   me.TopXDirection,
		"topX":            me.TopX,
		"topXAggregation": me.TopXAggregation,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *DimensionDefinition) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"name":            &me.Name,
		"dimension":       &me.Dimension,
		"placeholders":    &me.Placeholders,
		"topXDirection":   &me.TopXDirection,
		"topX":            &me.TopX,
		"topXAggregation": &me.TopXAggregation,
	})
}

// TopXDirection How to calculate the **topX** values.
type TopXDirection string

// TopXDirections offers the known enum values
var TopXDirections = struct {
	Ascending  TopXDirection
	Descending TopXDirection
}{
	"ASCENDING",
	"DESCENDING",
}

// TopXAggregation The aggregation of the dimension.
type TopXAggregation string

// TopXAggregations offers the known enum values
var TopXAggregations = struct {
	Average         TopXAggregation
	Count           TopXAggregation
	Max             TopXAggregation
	Min             TopXAggregation
	OfInterestRatio TopXAggregation
	OtherRatio      TopXAggregation
	SingleValue     TopXAggregation
	Sum             TopXAggregation
}{
	"AVERAGE",
	"COUNT",
	"MAX",
	"MIN",
	"OF_INTEREST_RATIO",
	"OTHER_RATIO",
	"SINGLE_VALUE",
	"SUM",
}
