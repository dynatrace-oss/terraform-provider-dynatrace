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

package maintenance

import (
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MaintenanceWindow TODO: documentation
type MaintenanceWindow struct {
	Enabled           bool               `json:"enabled"`           // The maintenance window is enabled or disabled
	GeneralProperties *GeneralProperties `json:"generalProperties"` // The general properties of the maintenance window
	Schedule          *Schedule          `json:"schedule"`          // The schedule of the maintenance window
	Filters           Filters            `json:"filters,omitempty"` // The filters of the maintenance window
	LegacyID          *string            `json:"-"`
}

func (me *MaintenanceWindow) Name() string {
	return me.GeneralProperties.Name
}

func (me *MaintenanceWindow) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The maintenance window is enabled or disabled",
			Default:     true,
			Optional:    true,
		},
		"general_properties": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "The general properties of the maintenance window",
			Elem: &schema.Resource{
				Schema: new(GeneralProperties).Schema(),
			},
		},
		"schedule": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "The schedule of the maintenance window",
			Elem: &schema.Resource{
				Schema: new(Schedule).Schema(),
			},
		},
		"filters": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The filters of the maintenance window",
			Elem: &schema.Resource{
				Schema: new(Filters).Schema(),
			},
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of this setting when referred to by the Config REST API V1",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *MaintenanceWindow) EnsurePredictableOrder() {
	if len(me.Filters) > 0 {
		conds := []*Filter{}
		condStrings := sort.StringSlice{}
		for _, entry := range me.Filters {
			condBytes, _ := json.Marshal(entry)
			condStrings = append(condStrings, string(condBytes))
		}
		condStrings.Sort()
		for _, condString := range condStrings {
			cond := Filter{}
			json.Unmarshal([]byte(condString), &cond)
			conds = append(conds, &cond)
		}
		me.Filters = conds
	}
}

func (me *MaintenanceWindow) MarshalHCL(properties hcl.Properties) error {
	me.EnsurePredictableOrder()

	return properties.EncodeAll(map[string]any{
		"enabled":            me.Enabled,
		"general_properties": me.GeneralProperties,
		"schedule":           me.Schedule,
		"filters":            me.Filters,
		"legacy_id":          me.LegacyID,
	})
}

func (me *MaintenanceWindow) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":            &me.Enabled,
		"general_properties": &me.GeneralProperties,
		"schedule":           &me.Schedule,
		"filters":            &me.Filters,
		"legacy_id":          &me.LegacyID,
	})
}
