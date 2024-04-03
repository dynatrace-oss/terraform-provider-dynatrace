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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type SourceFilters []*SourceFilter

func (me *SourceFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SourceFilter).Schema()},
		},
	}
}

func (me SourceFilters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("source", me)
}

func (me *SourceFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("source", me)
}

// Ingest source filter. The source filter is matched against the source of the ingest data. This way a subset of a specified data source can be used for creating the type.
type SourceFilter struct {
	Condition  *string          `json:"condition,omitempty"` // Specify a filter that needs to match in order for the extraction to happen.. Three different filters are supported: `$eq(value)` will ensure that the source matches exactly 'value', `$prefix(value)` will ensure that the source begins with exactly 'value', '$exists()' will ensure that any source with matching dimension filter exists.\nIf your value contains the characters '(', ')' or '\\~', you need to escape them by adding a '\\~' in front of them.
	SourceType IngestDataSource `json:"sourceType"`          // Possible Values: `BusinessEvents`, `Entities`, `Events`, `Logs`, `Metrics`, `Spans`, `Topology`
}

func (me *SourceFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "Specify a filter that needs to match in order for the extraction to happen.. Three different filters are supported: `$eq(value)` will ensure that the source matches exactly 'value', `$prefix(value)` will ensure that the source begins with exactly 'value', '$exists()' will ensure that any source with matching dimension filter exists.\nIf your value contains the characters '(', ')' or '\\~', you need to escape them by adding a '\\~' in front of them.",
			Optional:    true, // precondition
		},
		"source_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BusinessEvents`, `Entities`, `Events`, `Logs`, `Metrics`, `Spans`, `Topology`",
			Required:    true,
		},
	}
}

func (me *SourceFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":   me.Condition,
		"source_type": me.SourceType,
	})
}

func (me *SourceFilter) HandlePreconditions() error {
	if me.Condition == nil && !slices.Contains([]string{"Logs", "Spans", "Topology"}, string(me.SourceType)) {
		return fmt.Errorf("'condition' must be specified if 'source_type' is set to '%v'", me.SourceType)
	}
	return nil
}

func (me *SourceFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":   &me.Condition,
		"source_type": &me.SourceType,
	})
}
