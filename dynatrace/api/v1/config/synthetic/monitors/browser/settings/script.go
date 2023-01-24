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

package browser

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings/event"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Script struct {
	Version       string        `json:"version"`                 // Script version—use the `1.0` value here
	Type          ScriptType    `json:"type"`                    // The type of monitor. Possible values are `clickpath` for clickpath monitors and `availability` for single-URL browser monitors. These monitors are only allowed to have one event of the `navigate` type
	Configuration *ScriptConfig `json:"configuration,omitempty"` // The setup of the monitor
	Events        event.Events  `json:"events,omitempty"`        // Steps of the clickpath—the first step must always be of the `navigate` type
}

func (me *Script) GetVersion() string {
	return "1.0"
}

func (me *Script) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of monitor. Possible values are `clickpath` for clickpath monitors and `availability` for single-URL browser monitors. These monitors are only allowed to have one event of the `navigate` type",
			Required:    true,
		},
		"configuration": {
			Type:        schema.TypeList,
			Description: "The setup of the monitor",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: new(ScriptConfig).Schema(),
			},
		},
		"events": {
			Type:        schema.TypeList,
			Description: "Steps of the clickpath—the first step must always be of the `navigate` type",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(event.Events).Schema()},
		},
	}
}

func (me *Script) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("configuration", me.Configuration); err != nil {
		return err
	}
	if err := properties.Encode("events", me.Events); err != nil {
		return err
	}
	return nil
}

func (me *Script) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Version = me.GetVersion()
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	if err := decoder.Decode("configuration", &me.Configuration); err != nil {
		return err
	}
	if err := decoder.Decode("events", &me.Events); err != nil {
		return err
	}
	return nil
}
