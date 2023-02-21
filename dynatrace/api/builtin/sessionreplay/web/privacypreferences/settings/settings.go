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

package privacypreferences

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID           *string              `json:"-" scope:"applicationId"` // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	EnableOptInMode         bool                 `json:"enableOptInMode"`         // (Field has overlap with `dynatrace_application_data_privacy`) When [Session Replay opt-in mode](https://dt-url.net/sr-opt-in-mode) is turned on, Session Replay is deactivated until explicitly activated via an API call.
	MaskingPresets          *MaskingPresetConfig `json:"maskingPresets"`          // (Field has overlap with `dynatrace_application_data_privacy`) To protect your end users' privacy, select or customize [predefined masking options](https://dt-url.net/sr-masking-preset-options) that suit your content recording and playback requirements.
	UrlExclusionPatternList []string             `json:"urlExclusionPatternList"` // (Field has overlap with `dynatrace_application_data_privacy`) Exclude webpages or views from Session Replay recording by adding [URL exclusion rules](https://dt-url.net/sr-url-exclusion)
}

func (me *Settings) Name() string {
	return *me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
		},
		"enable_opt_in_mode": {
			Type:        schema.TypeBool,
			Description: "(Field has overlap with `dynatrace_application_data_privacy`) When [Session Replay opt-in mode](https://dt-url.net/sr-opt-in-mode) is turned on, Session Replay is deactivated until explicitly activated via an API call.",
			Required:    true,
		},
		"masking_presets": {
			Type:        schema.TypeList,
			Description: "(Field has overlap with `dynatrace_application_data_privacy`) To protect your end users' privacy, select or customize [predefined masking options](https://dt-url.net/sr-masking-preset-options) that suit your content recording and playback requirements.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(MaskingPresetConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"url_exclusion_pattern_list": {
			Type:        schema.TypeSet,
			Description: "(Field has overlap with `dynatrace_application_data_privacy`) Exclude webpages or views from Session Replay recording by adding [URL exclusion rules](https://dt-url.net/sr-url-exclusion)",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":             me.ApplicationID,
		"enable_opt_in_mode":         me.EnableOptInMode,
		"masking_presets":            me.MaskingPresets,
		"url_exclusion_pattern_list": me.UrlExclusionPatternList,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":             &me.ApplicationID,
		"enable_opt_in_mode":         &me.EnableOptInMode,
		"masking_presets":            &me.MaskingPresets,
		"url_exclusion_pattern_list": &me.UrlExclusionPatternList,
	})
}
