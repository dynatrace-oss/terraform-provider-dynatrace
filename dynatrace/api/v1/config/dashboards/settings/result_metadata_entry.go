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

type ResultMetadataEntry struct {
	Key      string
	Config   *CustomChartingItemMetadataConfig
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *ResultMetadataEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A generated key by the Dynatrace Server",
		},
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

func (me *ResultMetadataEntry) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("last_modified", int(opt.Int64(me.Config.LastModified))); err != nil {
		return err
	}
	if err := properties.Encode("custom_color", me.Config.CustomColor); err != nil {
		return err
	}
	if err := properties.Encode("key", me.Key); err != nil {
		return err
	}
	return nil
}

func (me *ResultMetadataEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "key")
		delete(me.Unknowns, "last_modified")
		delete(me.Unknowns, "custom_color")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("key"); ok {
		me.Key = value.(string)
	}
	if value, ok := decoder.GetOk("last_modified"); ok {
		if me.Config == nil {
			me.Config = new(CustomChartingItemMetadataConfig)
		}
		me.Config.LastModified = opt.NewInt64(int64(value.(int)))
	}
	if value, ok := decoder.GetOk("custom_color"); ok {
		if me.Config == nil {
			me.Config = new(CustomChartingItemMetadataConfig)
		}
		me.Config.CustomColor = value.(string)
	}
	return nil
}
