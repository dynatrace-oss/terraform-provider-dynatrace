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

package usersettings

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AutoLanguage bool      `json:"auto-language"`      // Language - use browser default
	AutoRegion   bool      `json:"auto-region"`        // Region - use browser default
	AutoTheme    bool      `json:"auto-theme"`         // Theme - use browser default
	AutoTimezone bool      `json:"auto-timezone"`      // Timezone - use browser default
	Language     *Language `json:"language,omitempty"` // Possible Values: `En`, `Ja`
	Region       *string   `json:"region,omitempty"`   // Region
	Scope        string    `json:"-" scope:"scope"`    // The scope of this setting (user, userdefaults)
	Theme        *Theme    `json:"theme,omitempty"`    // Possible Values: `Dark`, `Light`
	Timezone     *string   `json:"timezone,omitempty"` // Timezone
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_language": {
			Type:        schema.TypeBool,
			Description: "Language - use browser default",
			Required:    true,
		},
		"auto_region": {
			Type:        schema.TypeBool,
			Description: "Region - use browser default",
			Required:    true,
		},
		"auto_theme": {
			Type:        schema.TypeBool,
			Description: "Theme - use browser default",
			Required:    true,
		},
		"auto_timezone": {
			Type:        schema.TypeBool,
			Description: "Timezone - use browser default",
			Required:    true,
		},
		"language": {
			Type:        schema.TypeString,
			Description: "Possible Values: `En`, `Ja`",
			Optional:    true, // precondition
		},
		"region": {
			Type:        schema.TypeString,
			Description: "Region",
			Optional:    true, // precondition
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (user, userdefaults)",
			Required:    true,
			ForceNew:    true,
		},
		"theme": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Dark`, `Light`",
			Optional:    true, // precondition
		},
		"timezone": {
			Type:        schema.TypeString,
			Description: "Timezone",
			Optional:    true, // precondition
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_language": me.AutoLanguage,
		"auto_region":   me.AutoRegion,
		"auto_theme":    me.AutoTheme,
		"auto_timezone": me.AutoTimezone,
		"language":      me.Language,
		"region":        me.Region,
		"scope":         me.Scope,
		"theme":         me.Theme,
		"timezone":      me.Timezone,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Timezone == nil) && (!me.AutoTimezone) {
		me.Timezone = opt.NewString("")
	}
	if (me.Language == nil) && (!me.AutoLanguage) {
		return fmt.Errorf("'language' must be specified if 'auto_language' is set to '%v'", me.AutoLanguage)
	}
	if (me.Language != nil) && (me.AutoLanguage) {
		return fmt.Errorf("'language' must not be specified if 'auto_language' is set to '%v'", me.AutoLanguage)
	}
	if (me.Region == nil) && (!me.AutoRegion) {
		return fmt.Errorf("'region' must be specified if 'auto_region' is set to '%v'", me.AutoRegion)
	}
	if (me.Theme == nil) && (!me.AutoTheme) {
		return fmt.Errorf("'theme' must be specified if 'auto_theme' is set to '%v'", me.AutoTheme)
	}
	if (me.Theme != nil) && (me.AutoTheme) {
		return fmt.Errorf("'theme' must not be specified if 'auto_theme' is set to '%v'", me.AutoTheme)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_language": &me.AutoLanguage,
		"auto_region":   &me.AutoRegion,
		"auto_theme":    &me.AutoTheme,
		"auto_timezone": &me.AutoTimezone,
		"language":      &me.Language,
		"region":        &me.Region,
		"scope":         &me.Scope,
		"theme":         &me.Theme,
		"timezone":      &me.Timezone,
	})
}
