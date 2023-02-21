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

// Content masking preferences. To protect your end users' privacy, select or customize [predefined masking options](https://dt-url.net/sr-masking-preset-options) that suit your content recording and playback requirements.// The recording masking settings are applied at record time to all webpages that your users navigate to. The playback masking settings are applied when replaying recorded sessions, including those that were recorded before your masking settings were applied.// Note: When you set the recording masking settings to a more restrictive option, the same option is also enabled for playback masking settings, which affects all past recorded sessions as well.
type MaskingPresetConfig struct {
	PlaybackMaskingAllowListRules  AllowListRules `json:"playbackMaskingAllowListRules,omitempty"`  // The elements are defined by the CSS selector or attribute name.
	PlaybackMaskingBlockListRules  BlockListRules `json:"playbackMaskingBlockListRules,omitempty"`  // The elements are defined by the CSS selector or attribute name.
	PlaybackMaskingPreset          MaskingPreset  `json:"playbackMaskingPreset"`                    // Possible Values: `MASK_ALL`, `MASK_USER_INPUT`, `ALLOW_LIST`, `BLOCK_LIST`
	RecordingMaskingAllowListRules AllowListRules `json:"recordingMaskingAllowListRules,omitempty"` // The elements are defined by the CSS selector or attribute name.
	RecordingMaskingBlockListRules BlockListRules `json:"recordingMaskingBlockListRules,omitempty"` // The elements are defined by the CSS selector or attribute name.
	RecordingMaskingPreset         MaskingPreset  `json:"recordingMaskingPreset"`                   // Possible Values: `MASK_USER_INPUT`, `ALLOW_LIST`, `BLOCK_LIST`, `MASK_ALL`
}

func (me *MaskingPresetConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"playback_masking_allow_list_rules": {
			Type:        schema.TypeList,
			Description: "The elements are defined by the CSS selector or attribute name.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(AllowListRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"playback_masking_block_list_rules": {
			Type:        schema.TypeList,
			Description: "The elements are defined by the CSS selector or attribute name.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(BlockListRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"playback_masking_preset": {
			Type:        schema.TypeString,
			Description: "Possible Values: `MASK_ALL`, `MASK_USER_INPUT`, `ALLOW_LIST`, `BLOCK_LIST`",
			Required:    true,
		},
		"recording_masking_allow_list_rules": {
			Type:        schema.TypeList,
			Description: "The elements are defined by the CSS selector or attribute name.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(AllowListRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"recording_masking_block_list_rules": {
			Type:        schema.TypeList,
			Description: "The elements are defined by the CSS selector or attribute name.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(BlockListRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"recording_masking_preset": {
			Type:        schema.TypeString,
			Description: "Possible Values: `MASK_USER_INPUT`, `ALLOW_LIST`, `BLOCK_LIST`, `MASK_ALL`",
			Required:    true,
		},
	}
}

func (me *MaskingPresetConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"playback_masking_allow_list_rules":  me.PlaybackMaskingAllowListRules,
		"playback_masking_block_list_rules":  me.PlaybackMaskingBlockListRules,
		"playback_masking_preset":            me.PlaybackMaskingPreset,
		"recording_masking_allow_list_rules": me.RecordingMaskingAllowListRules,
		"recording_masking_block_list_rules": me.RecordingMaskingBlockListRules,
		"recording_masking_preset":           me.RecordingMaskingPreset,
	})
}

func (me *MaskingPresetConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"playback_masking_allow_list_rules":  &me.PlaybackMaskingAllowListRules,
		"playback_masking_block_list_rules":  &me.PlaybackMaskingBlockListRules,
		"playback_masking_preset":            &me.PlaybackMaskingPreset,
		"recording_masking_allow_list_rules": &me.RecordingMaskingAllowListRules,
		"recording_masking_block_list_rules": &me.RecordingMaskingBlockListRules,
		"recording_masking_preset":           &me.RecordingMaskingPreset,
	})
}
