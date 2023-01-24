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

package comparisoninfo

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Number Comparison for `NUMBER` attributes.
type Number struct {
	BaseComparisonInfo
	Comparison NumberComparison `json:"comparison"`       // Operator of the comparision. You can reverse it by setting **negate** to `true`.
	Value      *float64         `json:"value,omitempty"`  // The value to compare to.
	Values     []float64        `json:"values,omitempty"` // The values to compare to.
}

func (me *Number) GetType() Type {
	return Types.Number
}

func (me *Number) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"values": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The values to compare to",
			Elem:        &schema.Schema{Type: schema.TypeFloat},
		},
		"value": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "The value to compare to",
		},
		"operator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Operator of the comparison. You can reverse it by setting `negate` to `true`. Possible values are `EQUALS`, `EQUALS_ANY_OF`, `EXISTS`, `GREATER_THAN`, `GREATER_THAN_OR_EQUAL`, `LOWER_THAN` and `LOWER_THAN_OR_EQUAL`",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Number) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.EncodeAll(map[string]any{
		"values":   me.Values,
		"value":    me.Value,
		"operator": me.Comparison,
		"unknowns": me.Unknowns,
	}); err != nil {
		return err
	}
	// if len(me.Values) > 0 {
	// 	if err := properties.Encode("values", me.Values); err != nil { return err }
	// }
	return nil
}

func (me *Number) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"values":   &me.Values,
		"value":    &me.Value,
		"operator": &me.Comparison,
		"unknowns": &me.Unknowns,
	})
}

func (me *Number) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"type":       me.GetType(),
		"negate":     me.Negate,
		"values":     me.Values,
		"value":      me.Value,
		"comparison": me.Comparison,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Number) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"negate":     &me.Negate,
		"values":     &me.Values,
		"value":      &me.Value,
		"comparison": &me.Comparison,
	})
}

// NumberComparison Operator of the comparision. You can reverse it by setting **negate** to `true`.
type NumberComparison string

// NumberComparisons offers the known enum values
var NumberComparisons = struct {
	Equals             NumberComparison
	EqualsAnyOf        NumberComparison
	Exists             NumberComparison
	GreaterThan        NumberComparison
	GreaterThanOrEqual NumberComparison
	LowerThan          NumberComparison
	LowerThanOrEqual   NumberComparison
}{
	"EQUALS",
	"EQUALS_ANY_OF",
	"EXISTS",
	"GREATER_THAN",
	"GREATER_THAN_OR_EQUAL",
	"LOWER_THAN",
	"LOWER_THAN_OR_EQUAL",
}
