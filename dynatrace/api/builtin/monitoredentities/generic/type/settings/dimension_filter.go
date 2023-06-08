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

package generictype

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DimensionFilters []*DimensionFilter

func (me *DimensionFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"required_dimension": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DimensionFilter).Schema()},
		},
	}
}

func (me DimensionFilters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("required_dimension", me)
}

func (me *DimensionFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("required_dimension", me)
}

// Ingest dimension filter. A dimension describes a property key which is present in the ingest data.
type DimensionFilter struct {
	Key          string  `json:"key"`                    // A dimension key which needs to exist in the ingest data to match this filter.
	ValuePattern *string `json:"valuePattern,omitempty"` // A dimension value pattern which needs to exist in the ingest data to match this filter.
}

func (me *DimensionFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "A dimension key which needs to exist in the ingest data to match this filter.",
			Required:    true,
		},
		"value_pattern": {
			Type:        schema.TypeString,
			Description: "A dimension value pattern which needs to exist in the ingest data to match this filter.",
			Optional:    true, // nullable
		},
	}
}

func (me *DimensionFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":           me.Key,
		"value_pattern": me.ValuePattern,
	})
}

func (me *DimensionFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":           &me.Key,
		"value_pattern": &me.ValuePattern,
	})
}
