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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type UserActionProperties []*UserActionProperty

func (me *UserActionProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"property": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "User Action Property",
			Elem:        &schema.Resource{Schema: new(UserActionProperty).Schema()},
		},
	}
}

func (me UserActionProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("property", me)
}

func (me *UserActionProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("property", me)
}

type UserActionProperty struct {
	Key       *string    `json:"key,omitempty"`       // The key of the action property we're checking.
	Value     *string    `json:"value,omitempty"`     // Only actions that have this value in the specified property are included in the metric calculation.
	From      *float64   `json:"from,omitempty"`      // Only actions that have a value greater than or equal to this are included in the metric calculation.
	To        *float64   `json:"to,omitempty"`        // Only actions that have a value less than or equal to this are included in the metric calculation.
	MatchType *MatchType `json:"matchType,omitempty"` // Specifies the match type of a string filter, e.g. using Contains or Equals.
}

func (me *UserActionProperty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "The key of the action property we're checking.",
			Optional:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Only actions that have this value in the specified property are included in the metric calculation.",
			Optional:    true,
		},
		"from": {
			Type:        schema.TypeFloat,
			Description: "Only actions that have a value greater than or equal to this are included in the metric calculation.",
			Optional:    true,
		},
		"to": {
			Type:        schema.TypeFloat,
			Description: "Only actions that have a value less than or equal to this are included in the metric calculation.",
			Optional:    true,
		},
		"match_type": {
			Type:        schema.TypeString,
			Description: "Specifies the match type of a string filter, e.g. using Contains or Equals.",
			Optional:    true,
		},
	}
}

func (me *UserActionProperty) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":        me.Key,
		"value":      me.Value,
		"from":       me.From,
		"to":         me.To,
		"match_type": me.MatchType,
	})
}

func (me *UserActionProperty) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":        &me.Key,
		"value":      &me.Value,
		"from":       &me.From,
		"to":         &me.To,
		"match_type": &me.MatchType,
	})
}

type MatchType string

var MatchTypes = struct {
	Contains MatchType
	Equals   MatchType
}{
	"Contains",
	"Equals",
}
