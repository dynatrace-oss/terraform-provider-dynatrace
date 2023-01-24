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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KeyStrokes struct {
	EventBase
	TextValue         string         `json:"textValue,omitempty"`  // The text to enter
	Masked            *bool          `json:"masked,omitempty"`     // Indicates whether the `textValue` is encrypted (`true`) or not (`false`)
	SimulateBlurEvent bool           `json:"simulateBlurEvent"`    // Defines whether to blur the text field when it loses focus.\nSet to `true` to trigger the blur the `textValue`
	Wait              *WaitCondition `json:"wait,omitempty"`       // The wait condition for the event—defines how long Dynatrace should wait before the next action is executed
	Validate          Validations    `json:"validate,omitempty"`   // The validation rule for the event—helps you verify that your browser monitor loads the expected page content or page element
	Target            *Target        `json:"target,omitempty"`     // The tab on which the page should open
	Credential        *Credential    `json:"credential,omitempty"` // Credentials for this event
}

type Credential struct {
	ID    string `json:"id"`
	Field string `json:"field"`
}

func (me *Credential) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vault_id": {
			Type:        schema.TypeString,
			Description: "The ID of the credential within the Credentials Vault",
			Required:    true,
		},
		"field": {
			Type:        schema.TypeString,
			Description: "Either `username` or `password`",
			Required:    true,
		},
	}
}

func (me *Credential) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"vault_id": me.ID,
		"field":    me.Field,
	})
}

func (me *Credential) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("vault_id", &me.ID); err != nil {
		return err
	}
	if err := decoder.Decode("field", &me.Field); err != nil {
		return err
	}
	return nil
}

func (me *KeyStrokes) GetType() Type {
	return Types.KeyStrokes
}

func (me *KeyStrokes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"credential": {
			Type:        schema.TypeList,
			Description: "Credentials for this event",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Credential).Schema()},
		},
		"text": {
			Type:        schema.TypeString,
			Description: "The text to enter. Must not be specified if `credentials` from the vault are being used",
			Optional:    true,
		},
		"masked": {
			Type:        schema.TypeBool,
			Description: "Indicates whether the `textValue` is encrypted (`true`) or not (`false`). Must not be specified if `credentials` from the vault are being used",
			Optional:    true,
		},
		"simulate_blur_event": {
			Type:        schema.TypeBool,
			Description: "Defines whether to blur the text field when it loses focus.\nSet to `true` to trigger the blur the `textValue`",
			Optional:    true,
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

func (me *KeyStrokes) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("target", me.Target); err != nil {
		return err
	}
	if err := properties.Encode("credential", me.Credential); err != nil {
		return err
	}
	if err := properties.Encode("wait", me.Wait); err != nil {
		return err
	}
	if err := properties.Encode("validate", me.Validate); err != nil {
		return err
	}
	if me.Credential == nil {
		if me.Masked != nil && *me.Masked {
			if err := properties.Encode("masked", me.Masked); err != nil {
				return err
			}
		} else {
			if err := properties.Encode("masked", false); err != nil {
				return err
			}
		}
	} else {
		if err := properties.Encode("masked", false); err != nil {
			return err
		}
	}
	if err := properties.Encode("text", me.TextValue); err != nil {
		return err
	}
	if err := properties.Encode("simulate_blur_event", me.SimulateBlurEvent); err != nil {
		return err
	}
	return nil
}

func (me *KeyStrokes) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Type = Types.Tap
	if err := decoder.Decode("text", &me.TextValue); err != nil {
		return err
	}
	if err := decoder.Decode("masked", &me.Masked); err != nil {
		return err
	}
	if err := decoder.Decode("simulate_blur_event", &me.SimulateBlurEvent); err != nil {
		return err
	}
	if err := decoder.Decode("wait", &me.Wait); err != nil {
		return err
	}
	if err := decoder.Decode("credential", &me.Credential); err != nil {
		return err
	}
	if me.Credential == nil && me.Masked == nil {
		me.Masked = opt.NewBool(false)
	}
	if err := decoder.Decode("validate", &me.Validate); err != nil {
		return err
	}
	if err := decoder.Decode("target", &me.Target); err != nil {
		return err
	}
	return nil
}
