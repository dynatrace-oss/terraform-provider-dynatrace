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

package ingestsources

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DavisEventProperties []*DavisEventProperty

func (me *DavisEventProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"property": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DavisEventProperty).Schema()},
		},
	}
}

func (me DavisEventProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("property", me)
}

func (me *DavisEventProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("property", me)
}

type DavisEventProperty struct {
	Key      string                   `json:"key"`
	Strategy *FieldExtractionStrategy `json:"strategy,omitempty"` // Strategy for field extraction. Possible values: `equals`, `startsWith`
	Value    *string                  `json:"value,omitempty"`
}

func (me *DavisEventProperty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Required:    true,
		},
		"strategy": {
			Type:        schema.TypeString,
			Description: "Strategy for field extraction. Possible values: `equals`, `startsWith`",
			Optional:    true, // nullable
		},
		"value": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true, // precondition
		},
	}
}

func (me *DavisEventProperty) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":      me.Key,
		"strategy": me.Strategy,
		"value":    me.Value,
	})
}

func (me *DavisEventProperty) HandlePreconditions() error {
	if (me.Value != nil) && ((me.Strategy == nil || string(*me.Strategy) != "equals") && (me.Strategy != nil)) {
		return fmt.Errorf("'value' must not be specified unless ('strategy' is set to 'equals' or 'strategy' is not set); got 'strategy'='%v'", opt.ValOrNil(me.Strategy))
	}
	if (me.Value == nil) && ((me.Strategy != nil && string(*me.Strategy) == "equals") || (me.Strategy == nil)) {
		return fmt.Errorf("'value' must be specified when ('strategy' is set to 'equals' or 'strategy' is not set); got 'strategy'='%v'", opt.ValOrNil(me.Strategy))
	}
	return nil
}

func (me *DavisEventProperty) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":      &me.Key,
		"strategy": &me.Strategy,
		"value":    &me.Value,
	})
}
