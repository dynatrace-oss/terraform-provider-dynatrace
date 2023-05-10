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

type SessionReplay struct {
	OnCrash bool `json:"onCrash"` // Before enabling, Dynatrace checks your system against the [prerequisites for Session Replay](https://dt-url.net/t23s0ppi).
}

func (me *SessionReplay) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"on_crash": {
			Type:        schema.TypeBool,
			Description: "Before enabling, Dynatrace checks your system against the [prerequisites for Session Replay](https://dt-url.net/t23s0ppi).",
			Required:    true,
		},
	}
}

func (me *SessionReplay) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"on_crash": me.OnCrash,
	})
}

func (me *SessionReplay) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"on_crash": &me.OnCrash,
	})
}
