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

package privacy

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DoNotTrack struct {
	ComplyWithDoNotTrack bool              `json:"complyWithDoNotTrack"` // Comply with \"Do Not Track\" browser settings
	DoNotTrack           *DoNotTrackOption `json:"doNotTrack,omitempty"` // Possible Values: `Anonymous`, `Disable_rum`
}

func (me *DoNotTrack) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"comply_with_do_not_track": {
			Type:        schema.TypeBool,
			Description: "Comply with \"Do Not Track\" browser settings",
			Required:    true,
		},
		"do_not_track": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Anonymous`, `Disable_rum`",
			Optional:    true, // precondition
		},
	}
}

func (me *DoNotTrack) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"comply_with_do_not_track": me.ComplyWithDoNotTrack,
		"do_not_track":             me.DoNotTrack,
	})
}

func (me *DoNotTrack) HandlePreconditions() error {
	if me.DoNotTrack == nil && me.ComplyWithDoNotTrack {
		return fmt.Errorf("'do_not_track' must be specified if 'comply_with_do_not_track' is set to '%v'", me.ComplyWithDoNotTrack)
	}
	if me.DoNotTrack != nil && !me.ComplyWithDoNotTrack {
		return fmt.Errorf("'do_not_track' must not be specified if 'comply_with_do_not_track' is set to '%v'", me.ComplyWithDoNotTrack)
	}
	return nil
}

func (me *DoNotTrack) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"comply_with_do_not_track": &me.ComplyWithDoNotTrack,
		"do_not_track":             &me.DoNotTrack,
	})
}
