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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SessionReplay struct {
	CostAndTrafficControl    *int  `json:"costAndTrafficControl,omitempty"`    // Percentage of user sessions recorded with Session Replay. For example, if you have 50% for RUM and 50% for Session Replay, it results in 25% of sessions recorded with Session Replay.
	FullSessionReplay        *bool `json:"fullSessionReplay,omitempty"`        // Before enabling, Dynatrace checks your system against the [prerequisites for Session Replay](https://dt-url.net/t23s0ppi).
	FullSessionReplayOnGrail *bool `json:"fullSessionReplayOnGrail,omitempty"` // Enable New Session Replay Experience
	OnCrash                  bool  `json:"onCrash"`                            // Capture screen recordings that replay the user actions preceding all detected crashes. Before enabling, Dynatrace checks your system against the [prerequisites for Session Replay](https://dt-url.net/t23s0ppi).
	OnCrashOnGrail           *bool `json:"onCrashOnGrail,omitempty"`           // Enable New Session Replay on Crashes Experience
}

func (me *SessionReplay) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cost_and_traffic_control": {
			Type:        schema.TypeInt,
			Description: "Percentage of user sessions recorded with Session Replay. For example, if you have 50% for RUM and 50% for Session Replay, it results in 25% of sessions recorded with Session Replay.",
			Optional:    true, // nullable
		},
		"full_session_replay": {
			Type:        schema.TypeBool,
			Description: "Before enabling, Dynatrace checks your system against the [prerequisites for Session Replay](https://dt-url.net/t23s0ppi).",
			Optional:    true, // nullable
		},
		"full_session_replay_on_grail": {
			Type:        schema.TypeBool,
			Description: "Enable New Session Replay Experience",
			Optional:    true, // nullable & precondition
		},
		"on_crash": {
			Type:        schema.TypeBool,
			Description: "Capture screen recordings that replay the user actions preceding all detected crashes. Before enabling, Dynatrace checks your system against the [prerequisites for Session Replay](https://dt-url.net/t23s0ppi).",
			Required:    true,
		},
		"on_crash_on_grail": {
			Type:        schema.TypeBool,
			Description: "Enable New Session Replay on Crashes Experience",
			Optional:    true, // nullable & precondition
		},
	}
}

func (me *SessionReplay) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cost_and_traffic_control":     me.CostAndTrafficControl,
		"full_session_replay":          me.FullSessionReplay,
		"full_session_replay_on_grail": me.FullSessionReplayOnGrail,
		"on_crash":                     me.OnCrash,
		"on_crash_on_grail":            me.OnCrashOnGrail,
	})
}

func (me *SessionReplay) HandlePreconditions() error {
	if (me.FullSessionReplayOnGrail == nil) && (me.FullSessionReplay != nil && *me.FullSessionReplay) {
		me.FullSessionReplayOnGrail = new(false)
	}
	if (me.OnCrashOnGrail == nil) && (me.OnCrash) {
		me.OnCrashOnGrail = new(false)
	}
	if (me.FullSessionReplayOnGrail != nil) && (me.FullSessionReplay == nil || !*me.FullSessionReplay) {
		return fmt.Errorf("'full_session_replay_on_grail' must not be specified unless 'full_session_replay' is set to 'true'; got 'full_session_replay'='%v'", opt.ValOrNil(me.FullSessionReplay))
	}
	if (me.OnCrashOnGrail != nil) && (!me.OnCrash) {
		return fmt.Errorf("'on_crash_on_grail' must not be specified unless 'on_crash' is set to 'true'; got 'on_crash'='%v'", me.OnCrash)
	}
	return nil
}

func (me *SessionReplay) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cost_and_traffic_control":     &me.CostAndTrafficControl,
		"full_session_replay":          &me.FullSessionReplay,
		"full_session_replay_on_grail": &me.FullSessionReplayOnGrail,
		"on_crash":                     &me.OnCrash,
		"on_crash_on_grail":            &me.OnCrashOnGrail,
	})
}
