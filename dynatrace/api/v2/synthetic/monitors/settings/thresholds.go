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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Thresholds []*Threshold

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"threshold": {
			Type:        schema.TypeSet,
			Description: "The list of performance threshold rules",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Threshold).Schema()},
		},
	}
}

func (me Thresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("threshold", me)
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("threshold", me)
}

type Threshold struct {
	Aggregation       *Aggregation `json:"aggregation,omitempty"`       // Aggregation type, possible values: `AVG`, `MAX`, `MIN`
	DealertingSamples *int32       `json:"dealertingSamples,omitempty"` // Number of most recent non-violating request executions that closes the problem
	Samples           *int32       `json:"samples,omitempty"`           // Number of request executions in analyzed sliding window (sliding window size)
	StepIndex         *int32       `json:"stepIndex,omitempty"`         // Specify the step's index to which a threshold applies
	Threshold         *int32       `json:"threshold,omitempty"`         // Notify if monitor request takes longer than X milliseconds to execute
	ViolatingSamples  *int32       `json:"violatingSamples,omitempty"`  // Number of violating request executions in analyzed sliding window
}

func (me *Threshold) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aggregation": {
			Type:        schema.TypeString,
			Description: "Aggregation type, possible values: `AVG`, `MAX`, `MIN`",
			Optional:    true,
		},
		"dealerting_samples": {
			Type:        schema.TypeInt,
			Description: "Number of most recent non-violating request executions that closes the problem",
			Optional:    true,
		},
		"samples": {
			Type:        schema.TypeInt,
			Description: "Number of request executions in analyzed sliding window (sliding window size)",
			Optional:    true,
		},
		"step_index": {
			Type:        schema.TypeInt,
			Description: "Specify the step's index to which a threshold applies",
			Optional:    true,
		},
		"threshold": {
			Type:        schema.TypeInt,
			Description: "Notify if monitor request takes longer than X milliseconds to execute",
			Optional:    true,
		},
		"violating_samples": {
			Type:        schema.TypeInt,
			Description: "Number of violating request executions in analyzed sliding window",
			Optional:    true,
		},
	}
}

func (me *Threshold) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"aggregation":        me.Aggregation,
		"dealerting_samples": me.DealertingSamples,
		"samples":            me.Samples,
		"step_index":         me.StepIndex,
		"threshold":          me.Threshold,
		"violating_samples":  me.ViolatingSamples,
	})
}

func (me *Threshold) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"aggregation":        &me.Aggregation,
		"dealerting_samples": &me.DealertingSamples,
		"samples":            &me.Samples,
		"step_index":         &me.StepIndex,
		"threshold":          &me.Threshold,
		"violating_samples":  &me.ViolatingSamples,
	})
}
