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

// Entity A filter for the metrics entity dimensions.
type Entity struct {
	BaseDimension
	NameFilter *Filter `json:"nameFilter"` // A filter for a string value based on the given operator.
}

func (me *Entity) GetType() FilterType {
	return FilterTypes.Entity
}

func (me *Entity) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The dimensions key on the metric",
		},
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "A filter for a string value based on the given operator",
			Elem:        &schema.Resource{Schema: new(Filter).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Entity) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"key":    me.Key,
		"name":   me.Name,
		"index":  me.Index,
		"filter": me.NameFilter,
	})
}

func (me *Entity) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "key")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "index")
		delete(me.Unknowns, "filterType")
		delete(me.Unknowns, "nameFilter")

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
	if _, ok := decoder.GetOk("filter.#"); ok {
		me.NameFilter = new(Filter)
		if err := me.NameFilter.UnmarshalHCL(hcl.NewDecoder(decoder, "filter", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *Entity) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"filterType": me.GetType(),
		"key":        me.Key,
		"name":       me.Name,
		"index":      me.Index,
		"nameFilter": me.NameFilter,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Entity) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"filterType": &me.FilterType,
		"key":        &me.Key,
		"name":       &me.Name,
		"index":      &me.Index,
		"nameFilter": &me.NameFilter,
	}); err != nil {
		return err
	}
	return nil
}
