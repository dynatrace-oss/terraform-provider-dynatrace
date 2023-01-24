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

package event

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WaitCondition struct {
	WaitFor               string      `json:"waitFor"`                         // The time to wait before the next event is triggered. Possible values are `page_complete` (wait for the page to load completely), `network` (wait for background network activity to complete), `next_action` (wait for the next action), `time` (wait for a specified periodof time) and `validation` (wait for a specific element to appear)
	Milliseconds          *int        `json:"milliseconds,omitempty"`          // The time to wait, in millisencods. The maximum allowed value is `60000`. Required for the type `time`, not applicable otherwise.
	TimeoutInMilliseconds *int        `json:"timeoutInMilliseconds,omitempty"` // The maximum amount of time to wait for a certain element to appear, in milliseconds—if exceeded, the action is marked as failed.\nThe maximum allowed value is 60000. Required for the type `validation`, not applicable otherwise.
	Validation            *Validation `json:"validation,omitempty"`            // The element to wait for. Required for the `validation` type, not applicable otherwise.
}

func (me *WaitCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"wait_for": {
			Type:        schema.TypeString,
			Description: "The time to wait before the next event is triggered. Possible values are `page_complete` (wait for the page to load completely), `network` (wait for background network activity to complete), `next_action` (wait for the next action), `time` (wait for a specified periodof time) and `validation` (wait for a specific element to appear)",
			Required:    true,
		},
		"milliseconds": {
			Type:        schema.TypeInt,
			Description: "The time to wait, in millisencods. The maximum allowed value is `60000`. Required for the type `time`, not applicable otherwise.",
			Optional:    true,
		},
		"timeout": {
			Type:        schema.TypeInt,
			Description: "he maximum amount of time to wait for a certain element to appear, in milliseconds—if exceeded, the action is marked as failed.\nThe maximum allowed value is 60000. Required for the type `validation`, not applicable otherwise..",
			Optional:    true,
		},
		"validation": {
			Type:        schema.TypeList,
			Description: "The elements to wait for. Required for the `validation` type, not applicable otherwise.",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Validation).Schema()},
		},
	}
}

func (me *WaitCondition) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("wait_for", me.WaitFor); err != nil {
		return err
	}
	if err := properties.Encode("milliseconds", me.Milliseconds); err != nil {
		return err
	}
	if err := properties.Encode("timeout", me.TimeoutInMilliseconds); err != nil {
		return err
	}
	if me.Validation != nil {
		marshalled := hcl.Properties{}
		if err := me.Validation.MarshalHCL(marshalled); err == nil {
			properties["validation"] = []any{marshalled}
		}
	}
	return nil
}

func (me *WaitCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("wait_for", &me.WaitFor); err != nil {
		return err
	}
	if err := decoder.Decode("milliseconds", &me.Milliseconds); err != nil {
		return err
	}
	if err := decoder.Decode("timeout", &me.TimeoutInMilliseconds); err != nil {
		return err
	}
	if err := decoder.Decode("validation", &me.Validation); err != nil {
		return err
	}
	return nil
}
