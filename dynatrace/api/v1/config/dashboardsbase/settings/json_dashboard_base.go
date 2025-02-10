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

	"github.com/google/uuid"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JSONDashboardBase struct {
	Name     string
	Contents string
}

func (me *JSONDashboardBase) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"contents": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "This attribute exists for backwards compatibility. You do not have to define it.",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return true },
			StateFunc:        func(val any) string { return "" },
		},
	}
}

func (me *JSONDashboardBase) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{"contents": me.Contents})
	return nil
}

func (me *JSONDashboardBase) UnmarshalHCL(decoder hcl.Decoder) error {
	v, ok := decoder.GetOk("contents")
	if ok {
		me.Contents = v.(string)
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
	me.Contents = string(data)

	return nil
}
