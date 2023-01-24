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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DynamicFilters Dashboard filter configuration of a dashboard
type DynamicFilters struct {
	Filters            []string                   `json:"filters,omitempty"`            // A set of all possible global dashboard filters that can be applied to a dashboard \n\nCurrently supported values are: \n\n\tOS_TYPE,\n\tSERVICE_TYPE,\n\tDEPLOYMENT_TYPE,\n\tAPPLICATION_INJECTION_TYPE,\n\tPAAS_VENDOR_TYPE,\n\tDATABASE_VENDOR,\n\tHOST_VIRTUALIZATION_TYPE,\n\tHOST_MONITORING_MODE,\n\tKUBERNETES_CLUSTER,\n\tRELATED_CLOUD_APPLICATION,\n\tRELATED_NAMESPACE,\n\tTAG_KEY:<tagname>
	TagSuggestionTypes []string                   `json:"tagSuggestionTypes,omitempty"` // A set of entities applied for tag filter suggestions. You can fetch the list of possible values with the [GET all entity types](https://dt-url.net/dw03s7h)request. \n\nOnly applicable if the **filters** set includes `TAG_KEY:<tagname>`
	Unknowns           map[string]json.RawMessage `json:"-"`
}

func (me *DynamicFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filters": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "A set of all possible global dashboard filters that can be applied to a dashboard \n\nCurrently supported values are: \n\n\tOS_TYPE,\n\tSERVICE_TYPE,\n\tDEPLOYMENT_TYPE,\n\tAPPLICATION_INJECTION_TYPE,\n\tPAAS_VENDOR_TYPE,\n\tDATABASE_VENDOR,\n\tHOST_VIRTUALIZATION_TYPE,\n\tHOST_MONITORING_MODE,\n\tKUBERNETES_CLUSTER,\n\tRELATED_CLOUD_APPLICATION,\n\tRELATED_NAMESPACE,\n\tTAG_KEY:<tagname>",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tag_suggestion_types": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "A set of entities applied for tag filter suggestions. You can fetch the list of possible values with the [GET all entity types](https://dt-url.net/dw03s7h)request. \n\nOnly applicable if the **filters** set includes `TAG_KEY:<tagname>`",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *DynamicFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "filters")
		delete(me.Unknowns, "tag_suggestion_types")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if err := decoder.Decode("filters", &me.Filters); err != nil {
		return err
	}
	if err := decoder.Decode("tag_suggestion_types", &me.TagSuggestionTypes); err != nil {
		return err
	}
	return nil
}

func (me *DynamicFilters) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("filters", me.Filters); err != nil {
		return err
	}
	if err := properties.Encode("tag_suggestion_types", me.TagSuggestionTypes); err != nil {
		return err
	}
	return nil
}

func (me *DynamicFilters) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	m.Marshal("filters", me.Filters)
	m.Marshal("tagSuggestionTypes", me.TagSuggestionTypes)
	return json.Marshal(m)
}

func (me *DynamicFilters) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("filters", &me.Filters); err != nil {
		return err
	}
	if err := m.Unmarshal("tagSuggestionTypes", &me.TagSuggestionTypes); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
