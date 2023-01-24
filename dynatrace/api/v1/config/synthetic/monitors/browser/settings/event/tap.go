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

type Tap struct {
	EventBase
	Button   int            `json:"button"`             // the mouse button to be used for the click
	Wait     *WaitCondition `json:"wait,omitempty"`     // The wait condition for the event—defines how long Dynatrace should wait before the next action is executed
	Validate Validations    `json:"validate,omitempty"` // The validation rule for the event—helps you verify that your browser monitor loads the expected page content or page element
	Target   *Target        `json:"target,omitempty"`   // The tab on which the page should open
}

func (me *Tap) GetType() Type {
	return Types.Tap
}

func (me *Tap) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"button": {
			Type:        schema.TypeInt,
			Description: "the mouse button to be used for the click",
			Required:    true,
		},
		"wait": {
			Type:        schema.TypeList,
			Description: "The wait condition for the event—defines how long Dynatrace should wait before the next action is executed",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(WaitCondition).Schema()},
		},
		"validate": {
			Type:        schema.TypeList,
			Description: "The validation rules for the event—helps you verify that your browser monitor loads the expected page content or page element",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Validations).Schema()},
		},
		"target": {
			Type:        schema.TypeList,
			Description: "The tab on which the page should open",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Target).Schema()},
		},
	}
}

func (me *Tap) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("button", me.Button); err != nil {
		return err
	}
	if err := properties.Encode("wait", me.Wait); err != nil {
		return err
	}
	if err := properties.Encode("validate", me.Validate); err != nil {
		return err
	}
	if err := properties.Encode("target", me.Target); err != nil {
		return err
	}
	return nil
}

func (me *Tap) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Type = Types.Tap
	if err := decoder.Decode("button", &me.Button); err != nil {
		return err
	}
	if err := decoder.Decode("wait", &me.Wait); err != nil {
		return err
	}
	if err := decoder.Decode("validate", &me.Validate); err != nil {
		return err
	}
	if err := decoder.Decode("target", &me.Target); err != nil {
		return err
	}
	return nil
}
