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

package dashboardsbase

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/meta"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// var JSON_DASHBOARD_BASE_PLUS = os.Getenv("DYNATRACE_JSON_DASHBOARD_BASE_PLUS") == "true"
var JSON_DASHBOARD_BASE_PLUS = true
var DYNATRACE_DASHBOARD_TESTS = len(os.Getenv("DYNATRACE_DASHBOARD_TESTS")) > 0

type JSONDashboardBase struct {
	Name     string
	Contents string
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
	if bc == "dashboardMetadata" {
		delete(m, "popularity")
	}
	if bc == "" {
		ensure(m, "tiles", []any{})
	}

}

func (me *JSONDashboardBase) EnrichRequireds() *JSONDashboardBase {
	m := map[string]any{}
	json.Unmarshal([]byte(me.Contents), &m)
	enrm(m, "", false)
	data, _ := json.Marshal(m)
	me.Contents = string(data)
	return me
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

func (me *JSONDashboardBase) Schema() map[string]*schema.Schema {
	if JSON_DASHBOARD_BASE_PLUS {
		return map[string]*schema.Schema{
			"contents": {
				Type:        schema.TypeString,
				Optional:    true,
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

	return map[string]*schema.Schema{
		"contents": {
			Type:             schema.TypeString,
			Computed:         true,
			Optional:         true,
			Description:      "This attribute exists for backwards compatibility. You do not have to define it.",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return len(old) != 0 },
		},
	}
}

var lifecycleIgnoreChanges = meta.LifeCycle{
	IgnoreChanges: []string{
		"contents",
	},
}

func (me *JSONDashboardBase) MarshalHCL(properties hcl.Properties) error {
	if JSON_DASHBOARD_BASE_PLUS {
		return properties.EncodeAll(map[string]any{
			"contents":  me.Contents,
			"lifecycle": &lifecycleIgnoreChanges,
		})
	}
	return nil
}

func (me *JSONDashboardBase) UnmarshalHCL(decoder hcl.Decoder) error {
	if JSON_DASHBOARD_BASE_PLUS {
		v, ok := decoder.GetOk("contents")
		if ok {
			me.Contents = v.(string)
		}
	}
	return nil
}

func (me *JSONDashboardBase) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{"dashboardMetadata": map[string]any{"name": uuid.New().String(), "owner": uuid.New().String()}, "tiles": []string{}})
}

func (me *JSONDashboardBase) UnmarshalJSON(data []byte) error {
	reduced := struct {
		Metadata *DashboardMetadataBase `json:"dashboardMetadata"`
	}{}

	if err := json.Unmarshal(data, &reduced); err != nil {
		return err
	}

	var err error

	if data, err = json.Marshal(reduced); err != nil {
		return err
	}

	md := struct {
		Metadata *DashboardMetadataBase `json:"dashboardMetadata"`
	}{}
	err = json.Unmarshal(data, &md)
	if err != nil {
		return err
	}
	me.Name = md.Metadata.Name

	if JSON_DASHBOARD_BASE_PLUS {
		me.Contents = string(data)
	}

	return nil
}
