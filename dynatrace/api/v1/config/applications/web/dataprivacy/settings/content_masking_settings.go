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

package dataprivacy

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ContentMaskingSettings represents content masking settings for Session Replay. \n\nFor more details, see [Configure Session Replay](https://dt-url.net/0m03slq) in Dynatrace Documentation
type ContentMaskingSettings struct {
	RecordingMaskingSettingsVersion int32                        `json:"recordingMaskingSettingsVersion"` // The version of the content masking. \n\nYou can use this API only with the version 2. \n\nIf you're using version 1, set this field to `2` in the PUT request to switch to version 2
	RecordingMaskingSettings        *SessionReplayMaskingSetting `json:"recordingMaskingSettings"`        // Configuration of the Session Replay masking during Recording
	PlaybackMaskingSettings         *SessionReplayMaskingSetting `json:"playbackMaskingSettings"`         // Configuration of the Session Replay masking during Playback
}

func (me *ContentMaskingSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"recording": {
			Type:        schema.TypeList,
			Description: "Configuration of the Session Replay masking during Recording",
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SessionReplayMaskingSetting).Schema()},
		},
		"playback": {
			Type:        schema.TypeList,
			Description: "Configuration of the Session Replay masking during Playback",
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SessionReplayMaskingSetting).Schema()},
		},
	}
}

func (me *ContentMaskingSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"recording": me.RecordingMaskingSettings,
		"playback":  me.PlaybackMaskingSettings,
	})
}

func (me *ContentMaskingSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	me.RecordingMaskingSettingsVersion = 2
	return decoder.DecodeAll(map[string]any{
		"recording": &me.RecordingMaskingSettings,
		"playback":  &me.PlaybackMaskingSettings,
	})
}
