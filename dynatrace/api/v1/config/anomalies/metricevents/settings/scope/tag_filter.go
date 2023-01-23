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

package scope

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/common"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TagFilter A scope filter for tags on entities.
type TagFilter struct {
	BaseAlertingScope
	TagFilter *common.TagFilter `json:"tagFilter"` // A tag-based filter of monitored entities.
}

func (me *TagFilter) GetType() FilterType {
	return FilterTypes.Tag
}

func (me *TagFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "A filter for a string value based on the given operator",
			Elem:        &schema.Resource{Schema: new(common.TagFilter).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *TagFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.Encode("filter", me.TagFilter)
}

func (me *TagFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "filterType")
		delete(me.Unknowns, "tagFilter")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if _, ok := decoder.GetOk("filter.#"); ok {
		me.TagFilter = new(common.TagFilter)
		if err := me.TagFilter.UnmarshalHCL(hcl.NewDecoder(decoder, "filter", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *TagFilter) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"filterType": me.GetType(),
		"tagFilter":  me.TagFilter,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *TagFilter) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"filterType": &me.FilterType,
		"tagFilter":  &me.TagFilter,
	}); err != nil {
		return err
	}
	return nil
}
