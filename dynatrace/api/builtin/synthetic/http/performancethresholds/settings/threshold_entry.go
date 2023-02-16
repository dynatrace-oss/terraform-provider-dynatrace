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
	Event     string  `json:"event"`     // Request
	Threshold float64 `json:"threshold"` // Threshold (in seconds)
}

func (me *ThresholdEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event": {
			Type:        schema.TypeString,
			Description: "Request",
			Required:    true,
		},
		"threshold": {
			Type:        schema.TypeFloat,
			Description: "Threshold (in seconds)",
			Required:    true,
		},
	}
}

func (me *ThresholdEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event":     me.Event,
		"threshold": me.Threshold,
	})
}

func (me *ThresholdEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event":     &me.Event,
		"threshold": &me.Threshold,
	})
}
