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

package defaultversion

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DefaultVersion string  `json:"defaultVersion"`     // Default version
	Revision       *string `json:"revision,omitempty"` // Revision
}

func (me *Settings) Name() string {
	return "default_version"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_version": {
			Type:             schema.TypeString,
			Description:      "Default version",
			Required:         true,
			StateFunc:        func(i any) string { return "" },
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return true },
		},
		"revision": {
			Type:             schema.TypeString,
			Description:      "Revision",
			Optional:         true, // precondition
			StateFunc:        func(i any) string { return "" },
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return true },
		},
	}
}

func (me *Settings) Deprecated() string {
	return "The resource `dynatrace_oneagent_default_version` just exists for backwards compatibility. The settings it addresses are no longer configurable. The resource `dynatrace_oneagent_updates` covers this functionality now"
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"default_version": me.DefaultVersion,
		"revision":        me.Revision,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"default_version": &me.DefaultVersion,
		"revision":        &me.Revision,
	})
}
