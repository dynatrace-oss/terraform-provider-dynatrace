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

package useraction

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings/match"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Details Configuration of a user action-based conversion goal
type Details struct {
	Value         *string      `json:"value,omitempty"`       // The value to be matched to hit the conversion goal
	CaseSensitive bool         `json:"caseSensitive"`         // The match is case-sensitive (`true`) or (`false`)
	MatchType     *match.Type  `json:"matchType,omitempty"`   // The operator of the match. Possible values are `Begins`, `Contains` and `Ends`.
	MatchEntity   *MatchEntity `json:"matchEntity,omitempty"` // The type of the entity to which the rule applies. Possible values are `ActionName`, `CssSelector`, `JavaScriptVariable`, `MetaTag`, `PagePath`, `PageTitle`, `PageUrl`, `UrlAnchor` and `XhrUrl`.
	ActionType    *ActionType  `json:"actionType,omitempty"`  // Type of the action to which the rule applies. Possible values are `Custom`, `Load` and `Xhr`.
}

func (me *Details) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:        schema.TypeString,
			Description: "The value to be matched to hit the conversion goal",
			Optional:    true,
		},
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "The match is case-sensitive (`true`) or (`false`)",
			Optional:    true,
		},
		"match_type": {
			Type:        schema.TypeString,
			Description: "The operator of the match. Possible values are `Begins`, `Contains` and `Ends`.",
			Optional:    true,
		},
		"match_entity": {
			Type:        schema.TypeString,
			Description: "The type of the entity to which the rule applies. Possible values are `ActionName`, `CssSelector`, `JavaScriptVariable`, `MetaTag`, `PagePath`, `PageTitle`, `PageUrl`, `UrlAnchor` and `XhrUrl`.",
			Optional:    true,
		},
		"action_type": {
			Type:        schema.TypeString,
			Description: "Type of the action to which the rule applies. Possible values are `Custom`, `Load` and `Xhr`.",
			Optional:    true,
		},
	}
}

func (me *Details) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"value":          me.Value,
		"case_sensitive": me.CaseSensitive,
		"match_type":     me.MatchType,
		"match_entity":   me.MatchEntity,
		"action_type":    me.ActionType,
	})
}

func (me *Details) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"value":          &me.Value,
		"case_sensitive": &me.CaseSensitive,
		"match_type":     &me.MatchType,
		"match_entity":   &me.MatchEntity,
		"action_type":    &me.ActionType,
	})
}
