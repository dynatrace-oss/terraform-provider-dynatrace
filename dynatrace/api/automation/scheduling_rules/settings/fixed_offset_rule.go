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

package scheduling_rules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FixedOffsetRule struct {
	Rule   string `json:"rule"`   // Refers to a scheduling rule for which to produce valid days with an offset
	Offset int    `json:"offset"` // Every day of the scheduling rule referred to with `rule` will be offset by this amount of days
}

func (me *FixedOffsetRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeString,
			Description: "Refers to a scheduling rule for which to produce valid days with an offset",
			Required:    true,
		},
		"offset": {
			Type:             schema.TypeInt,
			Description:      "Every day of the scheduling rule referred to with `rule` will be offset by this amount of days",
			Required:         true,
			ValidateDiagFunc: ValidateRange(-50, 50),
		},
	}
}

func (me *FixedOffsetRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"rule":   me.Rule,
		"offset": me.Offset,
	})
}

func (me *FixedOffsetRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"rule":   &me.Rule,
		"offset": &me.Offset,
	})
}
