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

package sessionreplay

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SessionReplayDataPrivacySettings Data privacy settings for Session Replay
type SessionReplayDataPrivacySettings struct {
	OptIn                  bool                    `json:"optInModeEnabled,omitempty"`       // If `true`, session recording is disabled until JavaScriptAPI `dtrum.enableSessionReplay()` is called
	URLExclusionRules      []string                `json:"urlExclusionRules,omitempty"`      // A list of URLs to be excluded from recording
	ContentMaskingSettings *ContentMaskingSettings `json:"contentMaskingSettings,omitempty"` // Content masking settings for Session Replay. \n\nFor more details, see [Configure Session Replay](https://dt-url.net/0m03slq) in Dynatrace Documentation
}

func (me *SessionReplayDataPrivacySettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"opt_in": {
			Type:        schema.TypeBool,
			Description: "If `true`, session recording is disabled until JavaScriptAPI `dtrum.enableSessionReplay()` is called",
			Optional:    true,
		},
		"url_exclusion_rules": {
			Type:        schema.TypeList,
			Description: "A list of URLs to be excluded from recording",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"content_masking_settings": {
			Type:        schema.TypeList,
			Description: "Content masking settings for Session Replay. \n\nFor more details, see [Configure Session Replay](https://dt-url.net/0m03slq) in Dynatrace Documentation",
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ContentMaskingSettings).Schema()},
		},
	}
}

func (me *SessionReplayDataPrivacySettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"opt_in":                   me.OptIn,
		"url_exclusion_rules":      me.URLExclusionRules,
		"content_masking_settings": me.ContentMaskingSettings,
	})
}

func (me *SessionReplayDataPrivacySettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"opt_in":                   &me.OptIn,
		"url_exclusion_rules":      &me.URLExclusionRules,
		"content_masking_settings": &me.ContentMaskingSettings,
	})
}
