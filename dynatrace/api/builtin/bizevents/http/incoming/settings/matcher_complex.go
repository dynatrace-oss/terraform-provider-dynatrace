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

package incoming

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type MatcherComplexes []*MatcherComplex

func (me *MatcherComplexes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"trigger": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(MatcherComplex).Schema()},
		},
	}
}

func (me MatcherComplexes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("trigger", me)
}

func (me *MatcherComplexes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("trigger", me)
}

// Matcher. Rule must match
type MatcherComplex struct {
	CaseSensitive *bool              `json:"caseSensitive,omitempty"` // Case sensitive
	Source        *DataSourceComplex `json:"source"`
	Type          ComparisonEnum     `json:"type"` // Possible Values: `CONTAINS`, `ENDS_WITH`, `EQUALS`, `EXISTS`, `N_CONTAINS`, `N_ENDS_WITH`, `N_EQUALS`, `N_EXISTS`, `N_STARTS_WITH`, `STARTS_WITH`
	Value         *string            `json:"value,omitempty"`
}

func (me *MatcherComplex) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "Case sensitive",
			Optional:    true, // precondition
		},
		"source": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DataSourceComplex).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `ENDS_WITH`, `EQUALS`, `EXISTS`, `N_CONTAINS`, `N_ENDS_WITH`, `N_EQUALS`, `N_EXISTS`, `N_STARTS_WITH`, `STARTS_WITH`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
	}
}

func (me *MatcherComplex) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"case_sensitive": me.CaseSensitive,
		"source":         me.Source,
		"type":           me.Type,
		"value":          me.Value,
	})
}

func (me *MatcherComplex) HandlePreconditions() error {
	if me.CaseSensitive == nil && !slices.Contains([]string{"EXISTS", "N_EXISTS"}, string(me.Type)) {
		me.CaseSensitive = opt.NewBool(false)
	}
	if me.Value == nil && !slices.Contains([]string{"EXISTS", "N_EXISTS"}, string(me.Type)) {
		return fmt.Errorf("'value' must be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *MatcherComplex) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"case_sensitive": &me.CaseSensitive,
		"source":         &me.Source,
		"type":           &me.Type,
		"value":          &me.Value,
	})
}
