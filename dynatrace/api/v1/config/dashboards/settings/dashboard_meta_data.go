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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DashboardMetadata contains parameters of a dashboard
type DashboardMetadata struct {
	Name                string                     `json:"name"`                     // the name of the dashboard
	Shared              *bool                      `json:"shared,omitempty"`         // the dashboard is shared (`true`) or private (`false`)
	Owner               *string                    `json:"owner,omitempty"`          // the owner of the dashboard
	SharingDetails      *SharingInfo               `json:"sharingDetails,omitempty"` // represents sharing configuration of a dashboard
	Filter              *DashboardFilter           `json:"dashboardFilter,omitempty"`
	Tags                []string                   `json:"tags,omitempty"`                // a set of tags assigned to the dashboard
	Preset              bool                       `json:"preset"`                        // the dashboard is a preset (`true`)
	ValidFilterKeys     []string                   `json:"validFilterKeys,omitempty"`     // a set of all possible global dashboard filters that can be applied to dashboard
	DynamicFilters      *DynamicFilters            `json:"dynamicFilters,omitempty"`      // Dashboard filter configuration of a dashboard
	HasConsistentColors *bool                      `json:"hasConsistentColors,omitempty"` // the dashboard is a preset (`true`)
	TilesNameSize       *int                       `json:"tilesNameSize,omitempty"`       // no documentation available
	Unknowns            map[string]json.RawMessage `json:"-"`
}

func (me *DashboardMetadata) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the name of the dashboard",
			Required:    true,
		},
		"shared": {
			Type:        schema.TypeBool,
			Description: "the dashboard is shared (`true`) or private (`false`)",
			Optional:    true,
		},
		"consistent_colors": {
			Type:        schema.TypeBool,
			Description: "The tile uses consistent colors when rendering its content",
			Optional:    true,
		},
		"owner": {
			Type:        schema.TypeString,
			Description: "the owner of the dashboard",
			Required:    true,
		},
		"sharing_details": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "represents sharing configuration of a dashboard",
			Elem: &schema.Resource{
				Schema: new(SharingInfo).Schema(),
			},
		},
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Global filter Settings for the Dashboard",
			Elem: &schema.Resource{
				Schema: new(DashboardFilter).Schema(),
			},
		},
		"dynamic_filters": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Dashboard filter configuration of a dashboard",
			Elem: &schema.Resource{
				Schema: new(DynamicFilters).Schema(),
			},
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "a set of tags assigned to the dashboard",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"preset": {
			Type:        schema.TypeBool,
			Description: "the dashboard is a preset (`true`) or not (`false`). Default is `false`.",
			Optional:    true,
		},
		"tiles_name_size": {
			Type:        schema.TypeInt,
			Description: "No documentation available",
			Optional:    true,
		},
		"valid_filter_keys": {
			Type:        schema.TypeSet,
			Description: "a set of all possible global dashboard filters that can be applied to dashboard",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *DashboardMetadata) MarshalHCL(properties hcl.Properties) error {
	if len(me.Unknowns) > 0 {
		delete(me.Unknowns, "hasConsistentColors")
		delete(me.Unknowns, "tilesNameSize")
	}
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("shared", opt.Bool(me.Shared)); err != nil {
		return err
	}
	if err := properties.Encode("preset", me.Preset); err != nil {
		return err
	}
	if err := properties.Encode("consistent_colors", opt.Bool(me.HasConsistentColors)); err != nil {
		return err
	}
	if err := properties.Encode("owner", me.Owner); err != nil {
		return err
	}
	if err := properties.Encode("tags", me.Tags); err != nil {
		return err
	}
	if err := properties.Encode("valid_filter_keys", me.ValidFilterKeys); err != nil {
		return err
	}
	if err := properties.Encode("sharing_details", me.SharingDetails); err != nil {
		return err
	}
	if err := properties.Encode("filter", me.Filter); err != nil {
		return err
	}
	if err := properties.Encode("dynamic_filters", me.DynamicFilters); err != nil {
		return err
	}
	return nil
}

func (me *DashboardMetadata) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		if v, ok := me.Unknowns["tilesNameSize"]; ok {
			json.Unmarshal(v, &me.TilesNameSize)
		}
		delete(me.Unknowns, "tilesNameSize")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "shared")
		delete(me.Unknowns, "owner")
		delete(me.Unknowns, "sharing_details")
		delete(me.Unknowns, "dashboard_filter")
		delete(me.Unknowns, "dynamic_filters")
		delete(me.Unknowns, "tags")
		delete(me.Unknowns, "valid_filter_keys")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("tiles_name_size"); ok {
		me.TilesNameSize = opt.NewInt(value.(int))
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("shared"); ok {
		me.Shared = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("preset"); ok {
		me.Preset = value.(bool)
	}
	if value, ok := decoder.GetOk("owner"); ok {
		me.Owner = opt.NewString(value.(string))
	}
	if _, ok := decoder.GetOk("sharing_details.#"); ok {
		me.SharingDetails = new(SharingInfo)
		if err := me.SharingDetails.UnmarshalHCL(hcl.NewDecoder(decoder, "sharing_details", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("filter.#"); ok {
		me.Filter = new(DashboardFilter)
		if err := me.Filter.UnmarshalHCL(hcl.NewDecoder(decoder, "filter", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("dynamic_filters.#"); ok {
		me.DynamicFilters = new(DynamicFilters)
		if err := me.DynamicFilters.UnmarshalHCL(hcl.NewDecoder(decoder, "dynamic_filters", 0)); err != nil {
			return err
		}
	}
	if err := decoder.Decode("tags", &me.Tags); err != nil {
		return err
	}
	if value, ok := decoder.GetOk("preset"); ok {
		me.Preset = value.(bool)
	}
	if err := decoder.Decode("valid_filter_keys", &me.ValidFilterKeys); err != nil {
		return err
	}
	return nil
}

func (me *DashboardMetadata) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("name", me.Name); err != nil {
		return nil, err
	}
	if err := m.Marshal("tilesNameSize", me.TilesNameSize); err != nil {
		return nil, err
	}
	if err := m.Marshal("shared", me.Shared); err != nil {
		return nil, err
	}
	if err := m.Marshal("owner", me.Owner); err != nil {
		return nil, err
	}
	if err := m.Marshal("sharingDetails", me.SharingDetails); err != nil {
		return nil, err
	}
	if err := m.Marshal("dashboardFilter", me.Filter); err != nil {
		return nil, err
	}
	if err := m.Marshal("dynamicFilters", me.DynamicFilters); err != nil {
		return nil, err
	}
	if err := m.Marshal("tags", me.Tags); err != nil {
		return nil, err
	}
	if err := m.Marshal("preset", me.Preset); err != nil {
		return nil, err
	}
	if err := m.Marshal("validFilterKeys", me.ValidFilterKeys); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *DashboardMetadata) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("name", &me.Name); err != nil {
		return nil
	}
	if err := m.Unmarshal("tilesNameSize", &me.TilesNameSize); err != nil {
		return nil
	}
	if err := m.Unmarshal("shared", &me.Shared); err != nil {
		return nil
	}
	if err := m.Unmarshal("owner", &me.Owner); err != nil {
		return nil
	}
	if err := m.Unmarshal("sharingDetails", &me.SharingDetails); err != nil {
		return nil
	}
	if err := m.Unmarshal("dashboardFilter", &me.Filter); err != nil {
		return nil
	}
	if err := m.Unmarshal("dynamicFilters", &me.DynamicFilters); err != nil {
		return nil
	}
	if err := m.Unmarshal("tags", &me.Tags); err != nil {
		return nil
	}
	if err := m.Unmarshal("preset", &me.Preset); err != nil {
		return nil
	}
	if err := m.Unmarshal("validFilterKeys", &me.ValidFilterKeys); err != nil {
		return nil
	}
	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
