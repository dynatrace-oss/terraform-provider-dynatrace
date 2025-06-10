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

package enablement

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID       *string              `json:"-" scope:"applicationId"`       // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	ExperienceAnalytics *ExperienceAnalytics `json:"experienceAnalytics,omitempty"` // Experience Analytics
	Rum                 *Rum                 `json:"rum"`                           // Capture and analyze all user actions within your application. Enable [Real User Monitoring (RUM)](https://dt-url.net/1n2b0prq) to monitor and improve your application's performance, identify errors, and gain insight into your user's behavior and experience.
	SessionReplay       *SessionReplay       `json:"sessionReplay"`                 // [Session Replay](https://dt-url.net/session-replay) captures all user interactions within your application and replays them in a movie-like experience while providing [best-in-class security and data protection](https://dt-url.net/b303zxj).
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
			ForceNew:    true,
		},
		"experience_analytics": {
			Type:        schema.TypeList,
			Description: "Experience Analytics",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(ExperienceAnalytics).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rum": {
			Type:        schema.TypeList,
			Description: "Capture and analyze all user actions within your application. Enable [Real User Monitoring (RUM)](https://dt-url.net/1n2b0prq) to monitor and improve your application's performance, identify errors, and gain insight into your user's behavior and experience.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(Rum).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"session_replay": {
			Type:        schema.TypeList,
			Description: "[Session Replay](https://dt-url.net/session-replay) captures all user interactions within your application and replays them in a movie-like experience while providing [best-in-class security and data protection](https://dt-url.net/b303zxj).",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SessionReplay).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":       me.ApplicationID,
		"experience_analytics": me.ExperienceAnalytics,
		"rum":                  me.Rum,
		"session_replay":       me.SessionReplay,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":       &me.ApplicationID,
		"experience_analytics": &me.ExperienceAnalytics,
		"rum":                  &me.Rum,
		"session_replay":       &me.SessionReplay,
	})
}
