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

package synthetic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Dimensions []*Dimension

func (me *Dimensions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dimension": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A dimensions for the metric usage",
			Elem:        &schema.Resource{Schema: new(Dimension).Schema()},
		},
	}
}

func (me Dimensions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("dimension", me)
}

func (me *Dimensions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("dimension", me)
}

// Dimensions Parameters of a definition of a calculated synthetic metric.
type Dimension struct {
	TopX      *int32          `json:"topX,omitempty"` // The number of top values to be calculated.
	Dimension MetricDimension `json:"dimension"`      // The dimension of the metric. Possible values are `Event`, `Location`, `ResourceOrigin`
}

func (me *Dimension) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"top_x": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The number of top values to be calculated",
		},
		"dimension": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The dimension of the metric. Possible values are `Event`, `Location`, `ResourceOrigin`",
		},
	}
}

func (me *Dimension) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"top_x":     me.TopX,
		"dimension": me.Dimension,
	})
}

func (me *Dimension) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"top_x":     &me.TopX,
		"dimension": &me.Dimension,
	})
}

type MetricDimension string

var MetricDimensions = struct {
	Event          MetricDimension
	Location       MetricDimension
	ResourceOrigin MetricDimension
}{
	"Event",
	"Location",
	"ResourceOrigin",
}
