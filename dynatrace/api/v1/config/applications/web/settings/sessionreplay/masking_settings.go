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

package sessionreplay

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MaskingSetting represents configuration of the Session Replay masking
type MaskingSetting struct {
	Preset MaskingPreset `json:"maskingPreset"`          // The type of the masking: \n\n* `MASK_ALL`: Mask all texts, user input, and images. \n* `MASK_USER_INPUT`: Mask all data that is provided through user input \n* `ALLOW_LIST`: Only elements, specified in **maskingRules** are shown, everything else is masked. \n* `BLOCK_LIST`: Elements, specified in **maskingRules** are masked, everything else is shown.
	Rules  MaskingRules  `json:"maskingRules,omitempty"` // A list of masking rules
}

func (me *MaskingSetting) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"preset": {
			Type:        schema.TypeString,
			Description: "The type of the masking: \n\n* `MASK_ALL`: Mask all texts, user input, and images. \n* `MASK_USER_INPUT`: Mask all data that is provided through user input \n* `ALLOW_LIST`: Only elements, specified in **maskingRules** are shown, everything else is masked. \n* `BLOCK_LIST`: Elements, specified in **maskingRules** are masked, everything else is shown",
			Required:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "A list of masking rules",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(MaskingRules).Schema()},
		},
	}
}

func (me *MaskingSetting) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"preset": me.Preset,
		"rules":  me.Rules,
	})
}

func (me *MaskingSetting) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"preset": &me.Preset,
		"rules":  &me.Rules,
	})
}
