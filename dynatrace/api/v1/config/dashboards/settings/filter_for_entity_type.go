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

package dashboards

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FilterForEntityType struct {
	EntityType string
	Filters    []*FilterMatch
}

func (me *FilterForEntityType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_type": {
			Type:        schema.TypeString,
			Description: "The entity type (e.g. HOST, SERVICE, ...)",
			Required:    true,
		},
		"match": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "the tiles this Dashboard consist of",
			Elem: &schema.Resource{
				Schema: new(FilterMatch).Schema(),
			},
		},
	}
}

func (me *FilterForEntityType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("entity_type", me.EntityType); err != nil {
		return err
	}
	if err := properties.Encode("match", me.Filters); err != nil {
		return err
	}
	return nil
}

func (me *FilterForEntityType) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("entity_type"); ok {
		me.EntityType = value.(string)
	}
	if result, ok := decoder.GetOk("match.#"); ok {
		me.Filters = []*FilterMatch{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(FilterMatch)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "match", idx)); err != nil {
				return err
			}
			me.Filters = append(me.Filters, entry)
		}
	}
	return nil
}
