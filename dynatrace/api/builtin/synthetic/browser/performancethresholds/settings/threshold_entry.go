/**
* @license
* Copyright 2025 Dynatrace LLC
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

package performancethresholds

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ThresholdEntries []*ThresholdEntry

func (me *ThresholdEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"threshold": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ThresholdEntry).Schema()},
		},
	}
}

func (me ThresholdEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("threshold", me)
}

func (me *ThresholdEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("threshold", me)
}

type ThresholdEntry struct {
	DealertingSamples int     `json:"dealertingSamples"` // Number of most recent non-violating executions that closes the problem
	Event             string  `json:"event"`             // Synthetic event
	Samples           int     `json:"samples"`           // Number of executions in analyzed sliding window (sliding window size)
	Threshold         float64 `json:"threshold"`         // Threshold (in seconds)
	ViolatingSamples  int     `json:"violatingSamples"`  // Number of violating executions in analyzed sliding window
}

func (me *ThresholdEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dealerting_samples": {
			Type:        schema.TypeInt,
			Description: "Number of most recent non-violating executions that closes the problem",
			Optional:    true,
			Default:     5,
		},
		"event": {
			Type:        schema.TypeString,
			Description: "Synthetic event",
			Required:    true,
		},
		"samples": {
			Type:        schema.TypeInt,
			Description: "Number of executions in analyzed sliding window (sliding window size)",
			Optional:    true,
			Default:     5,
		},
		"threshold": {
			Type:        schema.TypeFloat,
			Description: "Threshold (in seconds)",
			Required:    true,
		},
		"violating_samples": {
			Type:        schema.TypeInt,
			Description: "Number of violating executions in analyzed sliding window",
			Optional:    true,
			Default:     3,
		},
	}
}

func (me *ThresholdEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dealerting_samples": me.DealertingSamples,
		"event":              me.Event,
		"samples":            me.Samples,
		"threshold":          me.Threshold,
		"violating_samples":  me.ViolatingSamples,
	})
}

func (me *ThresholdEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dealerting_samples": &me.DealertingSamples,
		"event":              &me.Event,
		"samples":            &me.Samples,
		"threshold":          &me.Threshold,
		"violating_samples":  &me.ViolatingSamples,
	})
}
