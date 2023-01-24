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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomChartingItemMetadataConfig Additional metadata for charted metric
type CustomChartingItemMetadataConfig struct {
	LastModified *int64                     `json:"lastModified,omitempty"` // The timestamp of the last metadata modification, in UTC milliseconds
	CustomColor  string                     `json:"customColor"`            // The color of the metric in the chart, hex format
	Unknowns     map[string]json.RawMessage `json:"-"`
}

func (me *CustomChartingItemMetadataConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"last_modified": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The timestamp of the last metadata modification, in UTC milliseconds",
		},
		"custom_color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The color of the metric in the chart, hex format",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomChartingItemMetadataConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("last_modified", int(opt.Int64(me.LastModified))); err != nil {
		return err
	}
	if err := properties.Encode("custom_color", me.CustomColor); err != nil {
		return err
	}
	return nil
}

func (me *CustomChartingItemMetadataConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "last_modified")
		delete(me.Unknowns, "custom_color")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("last_modified"); ok {
		me.LastModified = opt.NewInt64(int64(value.(int)))
	}
	if value, ok := decoder.GetOk("custom_color"); ok {
		me.CustomColor = value.(string)
	}
	return nil
}

func (me *CustomChartingItemMetadataConfig) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.LastModified != nil {
		rawMessage, err := json.Marshal(me.LastModified)
		if err != nil {
			return nil, err
		}
		m["lastModified"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.CustomColor)
		if err != nil {
			return nil, err
		}
		m["customColor"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *CustomChartingItemMetadataConfig) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["lastModified"]; found {
		if err := json.Unmarshal(v, &me.LastModified); err != nil {
			return err
		}
	}
	if v, found := m["customColor"]; found {
		if err := json.Unmarshal(v, &me.CustomColor); err != nil {
			return err
		}
	}
	delete(m, "lastModified")
	delete(m, "customColor")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
