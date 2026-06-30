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

package sitereliabilityguardian

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Segments []*Segment

func (me *Segments) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"segment": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Segment).Schema()},
		},
	}
}

func (me Segments) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("segment", me)
}

func (me *Segments) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("segment", me)
}

// Segment. A Grail segment reference that scopes a DQL objective query to a subset of observability data.
type Segment struct {
	ID        string           `json:"id"`                  // Dynatrace Grail segment ID that scopes the DQL query to data within the segment.
	Variables SegmentVariables `json:"variables,omitempty"` // Variables to parameterize the segment filter.
}

func (me *Segment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "Dynatrace Grail segment ID that scopes the DQL query to data within the segment.",
			Required:    true,
		},
		"variables": {
			Type:        schema.TypeList,
			Description: "Variables to parameterize the segment filter.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(SegmentVariables).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Segment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":        me.ID,
		"variables": me.Variables,
	})
}

func (me *Segment) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":        &me.ID,
		"variables": &me.Variables,
	})
}
