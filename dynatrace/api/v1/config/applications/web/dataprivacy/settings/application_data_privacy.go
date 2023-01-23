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

// ApplicationDataPrivacy represents data privacy settings of the application
type ApplicationDataPrivacy struct {
	WebApplicationID                *string                           `json:"-"`                                  // Dynatrace entity ID of the web application
	DataCaptureOptInEnabled         bool                              `json:"dataCaptureOptInEnabled"`            // Set to `true` to disable data capture and cookies until JavaScriptAPI `dtrum.enable()` is called
	PersistentCookieForUserTracking bool                              `json:"persistentCookieForUserTracking"`    // Set to `true` to set persistent cookie in order to recognize returning devices
	DoNotTrackBehaviour             DoNotTrackBehaviour               `json:"doNotTrackBehaviour"`                // How to handle the \"Do Not Track\" header: \n\n* `IGNORE_DO_NOT_TRACK`: ignore the header and capture the data. \n* `CAPTURE_ANONYMIZED`: capture the data but do not tie it to the user. \n* `DO_NOT_CAPTURE`: respect the header and do not capture.
	SessionReplayDataPrivacy        *SessionReplayDataPrivacySettings `json:"sessionReplayDataPrivacy,omitempty"` // Data privacy settings for Session Replay

	Name string `json:"-"`
}

func (me *ApplicationDataPrivacy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"web_application_id": {
			Type:        schema.TypeString,
			Description: "Dynatrace entity ID of the web application",
			Required:    true,
		},
		"data_capture_opt_in": {
			Type:        schema.TypeBool,
			Description: "Set to `true` to disable data capture and cookies until JavaScriptAPI `dtrum.enable()` is called",
			Optional:    true,
		},
		"persistent_cookie_for_user_tracking": {
			Type:        schema.TypeBool,
			Description: "Set to `true` to set persistent cookie in order to recognize returning devices",
			Optional:    true,
		},
		"do_not_track_behaviour": {
			Type:        schema.TypeString,
			Description: "How to handle the \"Do Not Track\" header: \n\n* `IGNORE_DO_NOT_TRACK`: ignore the header and capture the data. \n* `CAPTURE_ANONYMIZED`: capture the data but do not tie it to the user. \n* `DO_NOT_CAPTURE`: respect the header and do not capture.",
			Required:    true,
		},
		"session_replay_data_privacy": {
			Type:        schema.TypeList,
			Description: "Data privacy settings for Session Replay",
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SessionReplayDataPrivacySettings).Schema()},
		},
	}
}

func (me *ApplicationDataPrivacy) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"web_application_id":                  me.WebApplicationID,
		"data_capture_opt_in":                 me.DataCaptureOptInEnabled,
		"persistent_cookie_for_user_tracking": me.PersistentCookieForUserTracking,
		"do_not_track_behaviour":              me.DoNotTrackBehaviour,
		"session_replay_data_privacy":         me.SessionReplayDataPrivacy,
	})
}

func (me *ApplicationDataPrivacy) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"web_application_id":                  &me.WebApplicationID,
		"data_capture_opt_in":                 &me.DataCaptureOptInEnabled,
		"persistent_cookie_for_user_tracking": &me.PersistentCookieForUserTracking,
		"do_not_track_behaviour":              &me.DoNotTrackBehaviour,
		"session_replay_data_privacy":         &me.SessionReplayDataPrivacy,
	})
}
