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
	"os"
	"sort"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/assert"
	"golang.org/x/exp/slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JSONDashboard struct {
	name     string
	Contents string
}

var DYNATRACE_DASHBOARD_TESTS = len(os.Getenv("DYNATRACE_DASHBOARD_TESTS")) > 0

func denullslice(s []any) []any {
	if len(s) == 0 {
		return s
	}
	rs := []any{}
	for _, v := range s {
		switch tv := v.(type) {
		case map[string]any:
			m2 := denullmap(tv)
			if len(m2) > 0 {
				rs = append(rs, m2)
			}
		case []any:
			s2 := denullslice(tv)
			if len(s2) > 0 {
				rs = append(rs, s2)
			}
		default:
			rs = append(rs, v)
		}
	}
	return rs
}

func denullmap(m map[string]any) map[string]any {
	if len(m) == 0 {
		return m
	}
	for k, v := range m {
		switch tv := v.(type) {
		case map[string]any:
			m2 := denullmap(tv)
			if len(m2) == 0 {
				delete(m, k)
			} else {
				m[k] = m2
			}
		case []any:
			s2 := denullslice(tv)
			if len(s2) == 0 {
				delete(m, k)
			} else {
				m[k] = denullslice(tv)
			}
		case bool:
			if !tv {
				delete(m, k)
			}
		default:
			if v == nil {
				delete(m, k)
			}
		}
	}
	return m
}

func enrs(s []any, bc string, fordiff bool) {
	bc = strings.TrimPrefix(bc, ".")
	if len(s) == 0 {
		return
	}
	for _, v := range s {
		switch tv := v.(type) {
		case map[string]any:
			enrm(tv, bc+"[#]", fordiff)
		case []any:
			enrs(tv, bc+"[#]", fordiff)
		default:
		}
	}
}

func ensure(m map[string]any, key string, value any) {
	if _, found := m[key]; !found {
		m[key] = value
	}
}

func enrm(m map[string]any, bc string, fordiff bool) {
	bc = strings.TrimPrefix(bc, ".")
	if m == nil {
		return
	}
	if bc == "tiles[#].queries[#].filterBy" {
		if v, found := m["nestedFilters"]; found {
			if nestedFilters, ok := v.([]any); ok {
				jsons := []string{}
				for _, nestedFilter := range nestedFilters {
					if nestedFilterMap, ok := nestedFilter.(map[string]any); ok {
						if v, found := nestedFilterMap["criteria"]; found {
							if criteria, ok := v.([]any); ok {
								cjsons := []string{}
								for _, crit := range criteria {
									data, _ := json.Marshal(crit)
									cjsons = append(cjsons, string(data))
								}
								slices.Sort(cjsons)
								sortedCriteria := []any{}
								for _, js := range cjsons {
									nf := map[string]any{}
									json.Unmarshal([]byte(js), &nf)
									sortedCriteria = append(sortedCriteria, nf)
								}
								nestedFilterMap["criteria"] = sortedCriteria
								nestedFilter = nestedFilterMap
							}
						}
					}
					data, _ := json.Marshal(nestedFilter)
					jsons = append(jsons, string(data))
				}
				slices.Sort(jsons)
				sortedNestedFilters := []any{}
				for _, js := range jsons {
					nf := map[string]any{}
					json.Unmarshal([]byte(js), &nf)
					sortedNestedFilters = append(sortedNestedFilters, nf)
				}
				m["nestedFilters"] = sortedNestedFilters
			}
		}
	}
	if bc == "dashboardMetadata" {
		delete(m, "popularity")
	}
	if bc == "tiles[#].filterConfig" {
		ensure(m, "filtersPerEntityType", map[string]any{})
	}
	if bc == "tiles[#].filterConfig.chartConfig" {
		ensure(m, "resultMetadata", map[string]any{})
	}
	if bc == "tiles[#].filterConfig.chartConfig.series[#].dimensions[#]" {
		ensure(m, "resultMetadata", []any{})
	}
	if bc == "tiles[#].filterConfig.chartConfig" {
		ensure(m, "series", []any{})
		ensure(m, "legendShown", true)
	}
	if bc == "tiles[#].filterConfig.chartConfig.series[#]" {
		ensure(m, "dimensions", []any{})
	}
	if bc == "tiles[#].filterConfig.chartConfig.series[#].dimensions[#]" {
		ensure(m, "values", []string{})
		if DYNATRACE_DASHBOARD_TESTS && fordiff {
			delete(m, "name")
		}
	}
	if bc == "tiles[#].visualConfig" {
		ensure(m, "rules", []any{})
	}
	if bc == "tiles[#].queries[#]" {
		ensure(m, "enabled", true)
		ensure(m, "timeAggregation", "DEFAULT")
	}
	if bc == "" {
		ensure(m, "tiles", []any{})
	}
	if bc == "tiles[#]" {
		ensure(m, "queriesSettings", map[string]any{"resolution": ""})
		ensure(m, "assignedEntities", []any{})
		ensure(m, "query", "")
		ensure(m, "bounds", map[string]any{})
		// ensure(m, "visualConfig", map[string]any{})
		ensure(m, "visualConfig", nil)
		if name, found := m["name"]; found {
			ensure(m, "customName", name)
		}
		ensure(m, "type", "NOT_CONFIGURED")
		ensure(m, "visualizationConfig", map[string]any{})
		ensure(m, "isAutoRefreshDisabled", true)
	}
	if bc == "tiles[#].visualConfig" {
		ensure(m, "heatmapSettings", map[string]any{})
		ensure(m, "honeycombSettings", map[string]any{})
	}
	if bc == "tiles[#].visualConfig.heatmapSettings" {
		ensure(m, "yAxis", "VALUE")
	}
	if bc == "tiles[#].visualConfig.honeycombSettings" {
		ensure(m, "showHive", true)
		ensure(m, "showLegend", true)
	}
	if bc == "tiles[#].visualConfig.rules[#]" {
		ensure(m, "properties", map[string]any{"color": "DEFAULT"})
	}
	if bc == "tiles[#].visualConfig.rules[#].properties" {
		ensure(m, "color", "DEFAULT")
	}
	if bc == "tiles[#].visualConfig.thresholds[#]" {
		ensure(m, "visible", true)
	}
	if bc == "tiles[#].visualizationConfig" {
		ensure(m, "hasAxisBucketing", true)
	}
	if bc == "tiles[#].queries[#].filterBy.nestedFilters[#]" {
		if criteria, found := m["criteria"]; found {
			if crits, ok := criteria.([]any); ok {
				sortedCrits := []any{}
				marshCrits := []string{}
				for _, crit := range crits {
					data, _ := json.Marshal(crit)
					marshCrits = append(marshCrits, string(data))
				}
				slices.Sort(marshCrits)
				for _, marshCrit := range marshCrits {
					crit := map[string]any{}
					json.Unmarshal([]byte(marshCrit), &crit)
					sortedCrits = append(sortedCrits, crit)
				}
				m["criteria"] = sortedCrits
			}
		}
	}

	for k, v := range m {
		switch tv := v.(type) {
		case map[string]any:
			enrm(tv, bc+"."+k, fordiff)
		case []any:
			enrs(tv, bc+"."+k, fordiff)
		default:
		}
	}
	if bc == "tiles[#]" {
		if check(m, "tileType", "DTAQL", "PIE_CHART", "APPLICATION_WORLDMAP") {
			delete(m, "visualizationConfig")
		}
	}

}

func check(m map[string]any, key string, values ...string) bool {
	value, found := m[key]
	if !found {
		return false
	}
	svalue, ok := value.(string)
	if !ok {
		return false
	}
	if len(values) == 0 {
		return false
	}
	for _, v := range values {
		if v == svalue {
			return true
		}
	}
	return false
}

func (me *JSONDashboard) EnrichRequireds() *JSONDashboard {
	m := map[string]any{}
	json.Unmarshal([]byte(me.Contents), &m)
	enrm(m, "", false)
	data, _ := json.Marshal(m)
	me.Contents = string(data)
	return me
}

func (me *JSONDashboard) DeNull() {
	m := map[string]any{}
	json.Unmarshal([]byte(me.Contents), &m)
	denullmap(m)
	data, _ := json.Marshal(m)
	me.Contents = string(data)
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

func get(v any, key string) any {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]any); ok {
		if d, found := m[key]; found {
			return d
		}
		return nil
	}
	return nil
}

func diffSuppressedContent(content string) string {
	m := map[string]any{}
	json.Unmarshal([]byte(content), &m)
	if DYNATRACE_DASHBOARD_TESTS {
		if dmd := get(m, "dashboardMetadata"); dmd != nil {
			if df := get(dmd, "dashboardFilter"); df != nil {
				if mgmz := get(df, "managementZone"); mgmz != nil {
					delete(df.(map[string]any), "managementZone")
				}
			}
		}
	}
	denullmap(m)
	delete(m, "metadata")
	enrm(m, "", true)
	if DYNATRACE_DASHBOARD_TESTS {
		if tiles, found := m["tiles"]; found {
			if tileSlice, ok := tiles.([]any); ok {
				for _, tile := range tileSlice {
					if tm, ok := tile.(map[string]any); ok {
						delete(tm, "metricExpressions")
					}
				}

			}
		}
	}
	data, _ := json.Marshal(m)
	return string(data)
}

func (me *JSONDashboard) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"contents": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Contains the JSON Code of the Dashboard",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if len(old) == 0 && len(new) != 0 {
					return false
				}
				if len(old) != 0 && len(new) == 0 {
					return false
				}
				if !json.Valid([]byte(old)) || !json.Valid([]byte(new)) {
					return false
				}
				old = diffSuppressedContent(old)
				new = diffSuppressedContent(new)
				result := hcl.JSONStringsEqual(old, new)
				return result
			},
			StateFunc: func(val any) string {
				if json.Valid([]byte(val.(string))) {
					content := diffSuppressedContent(val.(string))
					return content
				}
				return val.(string)
			},
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
