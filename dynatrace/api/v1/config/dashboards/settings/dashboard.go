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

	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/common"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Dashboard the configuration of a dashboard
type Dashboard struct {
	ID       *string                    `json:"id,omitempty"`
	Metadata *DashboardMetadata         `json:"dashboardMetadata"` // contains parameters of a dashboard
	Tiles    []*Tile                    `json:"tiles"`             // the tiles the dashboard consists of
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *Dashboard) Name() string {
	return fmt.Sprintf("%s owned by %s", me.Metadata.Name, *me.Metadata.Owner)
}

func (me *Dashboard) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dashboard_metadata": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "contains parameters of a dashboard",
			Elem: &schema.Resource{
				Schema: new(DashboardMetadata).Schema(),
			},
		},
		"tile": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "the tiles this Dashboard consist of",
			Elem: &schema.Resource{
				Schema: new(Tile).Schema(),
			},
		},
		"metadata": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Description: "`metadata` exists for backwards compatibility but shouldn't get specified anymore",
			Deprecated:  "`metadata` exists for backwards compatibility but shouldn't get specified anymore",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(api.ConfigMetadata).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Dashboard) MarshalHCL(properties hcl.Properties) error {
	if len(me.Unknowns) > 0 {
		delete(me.Unknowns, "id")
		delete(me.Unknowns, "metadata")
		data, err := json.Marshal(me.Unknowns)
		if err != nil {
			return err
		}
		if err := properties.Encode("unknowns", string(data)); err != nil {
			return err
		}
	}
	if err := properties.Encode("dashboard_metadata", me.Metadata); err != nil {
		return err
	}
	if err := properties.Encode("tile", me.Tiles); err != nil {
		return err
	}
	return nil
}

func (me *Dashboard) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "dashboard_metadata")
		delete(me.Unknowns, "tile")
		delete(me.Unknowns, "id")
		delete(me.Unknowns, "metadata")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if _, ok := decoder.GetOk("dashboard_metadata.#"); ok {
		me.Metadata = new(DashboardMetadata)
		if err := me.Metadata.UnmarshalHCL(hcl.NewDecoder(decoder, "dashboard_metadata", 0)); err != nil {
			return err
		}
	}
	if result, ok := decoder.GetOk("tile.#"); ok {
		me.Tiles = []*Tile{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Tile)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "tile", idx)); err != nil {
				return err
			}
			me.Tiles = append(me.Tiles, entry)
		}
	}
	return nil
}

func (me *Dashboard) MarshalJSON() ([]byte, error) {
	if len(me.Unknowns) > 0 {
		delete(me.Unknowns, "metadata")
	}
	m := xjson.NewProperties(me.Unknowns)
	if me.ID != nil {
		if err := m.Marshal("id", *me.ID); err != nil {
			return nil, err
		}
	}
	if err := m.Marshal("dashboardMetadata", me.Metadata); err != nil {
		return nil, err
	}
	if err := m.Marshal("tiles", me.Tiles); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *Dashboard) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("id", &me.ID); err != nil {
		return err
	}
	if err := m.Unmarshal("dashboardMetadata", &me.Metadata); err != nil {
		return err
	}
	if err := m.Unmarshal("tiles", &me.Tiles); err != nil {
		return err
	}

	delete(m, "metadata")
	// delete(m, "id")
	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

func (me *Dashboard) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	delete(m, "id")
	return json.MarshalIndent(m, "", "  ")
}

func (me *Dashboard) Load(data []byte) error {
	var err error

	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "id")
	if data, err = json.Marshal(m); err != nil {
		return err
	}
	return json.Unmarshal(data, &me)
}

func (me *Dashboard) Anonymize() {
	me.ID = nil
}
