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

package requesterrors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CaptureSettings struct {
	Capture       bool  `json:"capture"`                 // Capture this error
	ConsiderForAI *bool `json:"considerForAi,omitempty"` // [View more details](https://dt-url.net/hd580p2k)
	ImpactApdex   *bool `json:"impactApdex,omitempty"`   // Include error in Apdex calculations
}

func (me *CaptureSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"capture": {
			Type:        schema.TypeBool,
			Description: "Capture this error",
			Required:    true,
		},
		"consider_for_ai": {
			Type:        schema.TypeBool,
			Description: "[View more details](https://dt-url.net/hd580p2k)",
			Optional:    true,
		},
		"impact_apdex": {
			Type:        schema.TypeBool,
			Description: "Include error in Apdex calculations",
			Optional:    true,
		},
	}
}

func (me *CaptureSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"capture":         me.Capture,
		"consider_for_ai": me.ConsiderForAI,
		"impact_apdex":    me.ImpactApdex,
	})
}

func (me *CaptureSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"capture":         &me.Capture,
		"consider_for_ai": &me.ConsiderForAI,
		"impact_apdex":    &me.ImpactApdex,
	})
	if me.Capture {
		if me.ConsiderForAI == nil {
			me.ConsiderForAI = opt.NewBool(false)
		}
		if me.ImpactApdex == nil {
			me.ImpactApdex = opt.NewBool(false)
		}
	}
	return err
}
