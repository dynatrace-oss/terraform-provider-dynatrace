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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DashboardMetadataBase contains parameters of a dashboard
type DashboardMetadataBase struct {
	Name  string  `json:"name"`            // the name of the dashboard
	Owner *string `json:"owner,omitempty"` // the owner of the dashboard
}

func (me *DashboardMetadataBase) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the name of the dashboard",
			Required:    true,
		},
		"owner": {
			Type:        schema.TypeString,
			Description: "the owner of the dashboard",
			Required:    true,
		},
	}
}

func (me *DashboardMetadataBase) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("owner", me.Owner); err != nil {
		return err
	}
	return nil
}

func (me *DashboardMetadataBase) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("owner"); ok {
		me.Owner = opt.NewString(value.(string))
	}
	return nil
}

func (me *DashboardMetadataBase) MarshalJSON() ([]byte, error) {
	m := xjson.Properties{}
	if err := m.Marshal("name", me.Name); err != nil {
		return nil, err
	}
	if err := m.Marshal("owner", me.Owner); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *DashboardMetadataBase) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("name", &me.Name); err != nil {
		return err
	}
	if err := m.Unmarshal("owner", &me.Owner); err != nil {
		return err
	}
	return nil
}
