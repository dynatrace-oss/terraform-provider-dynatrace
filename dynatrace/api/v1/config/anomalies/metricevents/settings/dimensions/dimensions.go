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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Dimensions []Dimension

func (me Dimensions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A filter for the metrics entity dimensions",
			Elem:        &schema.Resource{Schema: new(Entity).Schema()},
		},
		"string": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A filter for the metrics string dimensions",
			Elem:        &schema.Resource{Schema: new(String).Schema()},
		},
		"dimension": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A generic definition for a filter",
			Elem:        &schema.Resource{Schema: new(BaseDimension).Schema()},
		},
	}
}

func (me Dimensions) MarshalHCL(properties hcl.Properties) error {
	Entitys := []any{}
	Strings := []any{}
	baseDimensions := []map[string]any{}
	for _, dimension := range me {
		switch dim := dimension.(type) {
		case *Entity:
			marshalled := hcl.Properties{}
			if err := dim.MarshalHCL(marshalled); err == nil {
				Entitys = append(Entitys, marshalled)
			} else {
				return err
			}
		case *String:
			marshalled := hcl.Properties{}
			if err := dim.MarshalHCL(marshalled); err == nil {
				Strings = append(Strings, marshalled)
			} else {
				return err
			}
		case *BaseDimension:
			marshalled := hcl.Properties{}
			if err := dim.MarshalHCL(marshalled); err == nil {
				baseDimensions = append(baseDimensions, marshalled)
			} else {
				return err
			}
		default:
		}
	}
	if len(Entitys) > 0 {
		properties["entity"] = Entitys
	} else {
		properties["entity"] = nil

	}
	if len(Strings) > 0 {
		properties["string"] = Strings
	} else {
		properties["string"] = nil
	}
	if len(baseDimensions) > 0 {
		properties["dimension"] = baseDimensions
	} else {
		properties["dimension"] = nil
	}

	return nil
}

func (me *Dimensions) UnmarshalHCL(decoder hcl.Decoder) error {
	nme := Dimensions{}
	if result, ok := decoder.GetOk("entity.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Entity)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "entity", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("string.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(String)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "string", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("dimension.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(BaseDimension)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "dimension", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	*me = nme
	return nil
}

func (me *Dimensions) UnmarshalJSON(data []byte) error {
	dims := Dimensions{}
	rawMessages := []json.RawMessage{}
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}
	for _, rawMessage := range rawMessages {
		properties := map[string]json.RawMessage{}
		if err := json.Unmarshal(rawMessage, &properties); err != nil {
			return err
		}
		if rawFilterType, found := properties["filterType"]; found {
			var sFilterType string
			if err := json.Unmarshal(rawFilterType, &sFilterType); err != nil {
				return err
			}
			switch sFilterType {
			case string(FilterTypes.Entity):
				cfg := new(Entity)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.String):
				cfg := new(String)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			default:
				cfg := new(BaseDimension)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			}
		}
		*me = dims
	}
	return nil
}
