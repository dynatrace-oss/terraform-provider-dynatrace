/**
* @license
* Copyright 2026 Dynatrace LLC
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

package zones

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AlternativeZones []string     `json:"alternativeZones,omitempty"` // Network zones that should be used when the primary zone is not available.
	Description      *string      `json:"description,omitempty"`      // Description
	FallbackMode     FallbackMode `json:"fallbackMode"`               // Determines the network zone fallback behavior in case the primary and alternative zones are not available. Possible values: `ANY_ACTIVE_GATE`, `NONE`, `ONLY_DEFAULT_ZONE`
	ID               string       `json:"id"`                         // A lowercase string limited to 256 characters that can contain alphanumerics (0-9, a-z), hyphens (-), underscores (_), and dots (.), but can not start with a dot.
}

func (me *Settings) Name() string {
	return me.ID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alternative_zones": {
			Type:        schema.TypeSet,
			Description: "Network zones that should be used when the primary zone is not available.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"fallback_mode": {
			Type:        schema.TypeString,
			Description: "Determines the network zone fallback behavior in case the primary and alternative zones are not available. Possible values: `ANY_ACTIVE_GATE`, `NONE`, `ONLY_DEFAULT_ZONE`",
			Required:    true,
		},
		// renaming "id" to "identifier" because "id" is reserved.
		"identifier": {
			Type:        schema.TypeString,
			Description: "A lowercase string limited to 256 characters that can contain alphanumerics (0-9, a-z), hyphens (-), underscores (_), and dots (.), but can not start with a dot.",
			ForceNew:    true,
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"alternative_zones": me.AlternativeZones,
		"description":       me.Description,
		"fallback_mode":     me.FallbackMode,
		"identifier":        me.ID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"alternative_zones": &me.AlternativeZones,
		"description":       &me.Description,
		"fallback_mode":     &me.FallbackMode,
		"identifier":        &me.ID,
	})
}
