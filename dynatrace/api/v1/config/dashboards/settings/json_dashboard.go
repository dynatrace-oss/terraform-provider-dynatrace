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
	"fmt"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/assert"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JSONDashboard struct {
	name     string
	Contents string
}

func (me *JSONDashboard) Anonymize() {
	m := map[string]any{}
	json.Unmarshal([]byte(me.Contents), &m)
	delete(m, "id")
	data, _ := json.Marshal(m)
	me.Contents = string(data)
}

func (me *JSONDashboard) Equals(other any) (string, bool) {
	if o, ok := other.(*JSONDashboard); ok {
		ma := map[string]any{}
		json.Unmarshal([]byte(me.Contents), &ma)
		mb := map[string]any{}
		json.Unmarshal([]byte(o.Contents), &mb)
		return assert.Equals(ma, mb)
	}
	return "expected: JSONDashboard", false
}

func (me *JSONDashboard) Name() string {
	c := struct {
		Metadata *DashboardMetadata `json:"dashboardMetadata"`
	}{}
	json.Unmarshal([]byte(me.Contents), &c)

	if c.Metadata.Owner == nil {
		return fmt.Sprintf("%s owned by ", c.Metadata.Name)
	}
	return fmt.Sprintf("%s owned by %s", c.Metadata.Name, *c.Metadata.Owner)
}

func (me *JSONDashboard) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"contents": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Contains the JSON Code of the Dashboard",
		},
	}
}

func (me *JSONDashboard) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("contents", me.Contents)
}

func (me *JSONDashboard) UnmarshalHCL(decoder hcl.Decoder) error {
	v, ok := decoder.GetOk("contents")
	if ok {
		me.Contents = v.(string)
	}
	return nil
}

func (me *JSONDashboard) MarshalJSON() ([]byte, error) {
	return []byte(me.Contents), nil
}

func (me *JSONDashboard) UnmarshalJSON(data []byte) error {
	reduced := struct {
		Metadata *DashboardMetadata `json:"dashboardMetadata"`
		Tiles    []map[string]any   `json:"tiles"`
	}{}

	if err := json.Unmarshal(data, &reduced); err != nil {
		return err
	}

	var err error
	for _, tile := range reduced.Tiles {
		if v, found := tile["nameSize"]; found && v == nil {
			delete(tile, "nameSize")
		}
		if v, found := tile["tileFilter"]; found {
			vm := v.(map[string]any)
			if v, found := vm["managementZone"]; found && v == nil {
				delete(vm, "managementZone")
			}
			if v, found := vm["timeframe"]; found && v == nil {
				delete(vm, "timeframe")
			}
		}
		if v, found := tile["filterConfig"]; found && v != nil {
			if v, found := v.(map[string]any)["chartConfig"]; found {
				if v, found := v.(map[string]any)["series"]; found {
					for _, elem := range v.([]any) {
						if v, found := elem.(map[string]any)["percentile"]; found && v == nil {
							delete(elem.(map[string]any), "percentile")
						}
					}
				}
			}
		}
		var untypedQueries any
		var found bool
		if untypedQueries, found = tile["queries"]; found {
			queriesArr := untypedQueries.([]any)
			for _, untypedQuery := range queriesArr {
				query := untypedQuery.(map[string]any)
				var untypedFilterBy any
				if untypedFilterBy, found = query["filterBy"]; found && untypedFilterBy != nil {
					filterBy := untypedFilterBy.(map[string]any)
					var untypedNestedFilters any
					if untypedNestedFilters, found = filterBy["nestedFilters"]; found {
						nestedFilters := untypedNestedFilters.([]any)
						strs := []string{}
						for _, nestedFilter := range nestedFilters {
							var data []byte
							if data, err = json.Marshal(nestedFilter); err != nil {
								return err
							}
							strs = append(strs, string(data))
							sort.Strings(strs)
							nestedFilters = []any{}
							for _, str := range strs {
								nestedFilters = append(nestedFilters, json.RawMessage([]byte(str)))
							}
							filterBy["nestedFilters"] = nestedFilters
						}
					}

				}
			}
		}
	}

	if data, err = json.Marshal(reduced); err != nil {
		return err
	}

	md := struct {
		Metadata *DashboardMetadata `json:"dashboardMetadata"`
	}{}
	err = json.Unmarshal(data, &md)
	if err != nil {
		return err
	}
	me.name = md.Metadata.Name
	me.Contents = string(data)
	return nil
}
