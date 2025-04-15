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

package slo

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Criteria []*CriteriaDetail

func (me *Criteria) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"criteria_detail": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CriteriaDetail).Schema()},
		},
	}
}

func (me Criteria) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("criteria_detail", me)
}

func (me *Criteria) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("criteria_detail", me)
}

type CriteriaDetail struct {
	TimeframeFrom string   `json:"timeframeFrom" minlength:"3" maxlength:"30"`
	TimeframeTo   *string  `json:"timeframeTo,omitempty" minlength:"3" maxlength:"30"`
	Target        float32  `json:"target"`
	Warning       *float32 `json:"warning,omitempty"`
}

func (me *CriteriaDetail) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"timeframe_from": {
			Type:             schema.TypeString,
			Description:      "Timeframe from, example: `now-7d`",
			Required:         true,
			ValidateDiagFunc: Validate(ValidateMinLength(3), ValidateMaxLength(30)),
		},
		"timeframe_to": {
			Type:             schema.TypeString,
			Description:      "Timeframe to, example: `now`",
			Optional:         true,
			ValidateDiagFunc: Validate(ValidateMinLength(3), ValidateMaxLength(30)),
		},
		"target": {
			Type:        schema.TypeFloat,
			Description: "Criteria target, example: `99.8`",
			Required:    true,
		},
		"warning": {
			Type:        schema.TypeFloat,
			Description: "Criteria warning, example: `99.9`",
			Optional:    true,
		},
	}
}

func (me *CriteriaDetail) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"timeframe_from": me.TimeframeFrom,
		"timeframe_to":   me.TimeframeTo,
		"target":         me.Target,
		"warning":        me.Warning,
	})
}

func (me *CriteriaDetail) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"timeframe_from": &me.TimeframeFrom,
		"timeframe_to":   &me.TimeframeTo,
		"target":         &me.Target,
		"warning":        &me.Warning,
	})
}
