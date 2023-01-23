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

package dimensions

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MetricEventDimension A single filter for the metrics dimensions.
// This is the base version of the filter, depending on the type,
// the actual JSON may contain additional fields.
type Dimension interface {
	GetType() FilterType
}

// BaseDimension A single filter for the metrics dimensions.
// This is the base version of the filter, depending on the type,
// the actual JSON may contain additional fields.
type BaseDimension struct {
	FilterType FilterType                 `json:"filterType"`      // Defines the actual set of fields depending on the value. See one of the following objects:  * `ENTITY` -> MetricEventEntityDimensions  * `STRING` -> MetricEventStringDimensions
	Key        *string                    `json:"key,omitempty"`   // The dimensions key on the metric.
	Name       *string                    `json:"name,omitempty"`  // No documentation available
	Index      *int                       `json:"index,omitempty"` // No documentation available
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (me *BaseDimension) GetType() FilterType {
	return me.FilterType
}

func (me *BaseDimension) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The dimensions key on the metric",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "No documentation available",
		},
		"index": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "No documentation available",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the actual set of fields depending on the value",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *BaseDimension) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"key":   me.Key,
		"name":  me.Name,
		"index": me.Index,
		"type":  me.FilterType,
	})
}

func (me *BaseDimension) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "key")
		delete(me.Unknowns, "filterType")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "index")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("key"); ok {
		me.Key = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("index"); ok {
		me.Index = opt.NewInt(value.(int))
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.FilterType = FilterType(value.(string))
	}
	return nil
}

func (me *BaseDimension) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"filterType": me.FilterType,
		"key":        me.Key,
		"name":       me.Name,
		"index":      me.Index,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *BaseDimension) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"filterType": &me.FilterType,
		"key":        &me.Key,
		"name":       &me.Name,
		"index":      &me.Index,
	}); err != nil {
		return err
	}
	return nil
}

// DimensionFilterType Defines the actual set of fields depending on the value. See one of the following objects:
// * `ENTITY` -> MetricEventEntityDimensions
// * `STRING` -> MetricEventStringDimensions
type FilterType string

// FilterTypes offers the known enum values
var FilterTypes = struct {
	Entity FilterType
	String FilterType
}{
	"ENTITY",
	"STRING",
}
