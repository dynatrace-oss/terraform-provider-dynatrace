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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings/match"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DestinationDetails Configuration of a destination-based conversion goal
type DestinationDetails struct {
	URLOrPath     string      `json:"urlOrPath"`               // The path to be reached to hit the conversion goal
	MatchType     *match.Type `json:"matchType,omitempty"`     // The operator of the match. Possible values are `Begins`, `Contains` and `Ends`.
	CaseSensitive bool        `json:"caseSensitive,omitempty"` // The match is case-sensitive (`true`) or (`false`)
}

func (me *DestinationDetails) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url_or_path": {
			Type:        schema.TypeString,
			Description: "The path to be reached to hit the conversion goal",
			Required:    true,
		},
		"match_type": {
			Type:        schema.TypeString,
			Description: "The operator of the match. Possible values are `Begins`, `Contains` and `Ends`.",
			Optional:    true,
		},
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "The match is case-sensitive (`true`) or (`false`)",
			Optional:    true,
		},
	}
}

func (me *DestinationDetails) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"url_or_path":    me.URLOrPath,
		"match_type":     me.MatchType,
		"case_sensitive": me.CaseSensitive,
	})
}

func (me *DestinationDetails) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"url_or_path":    &me.URLOrPath,
		"match_type":     &me.MatchType,
		"case_sensitive": &me.CaseSensitive,
	})
}
