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

package dataminingblocklist

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataminingBlocklistEntries []*DataminingBlocklistEntry

func (me *DataminingBlocklistEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"blocklist_entrie": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DataminingBlocklistEntry).Schema()},
		},
	}
}

func (me DataminingBlocklistEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("blocklist_entrie", me)
}

func (me *DataminingBlocklistEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("blocklist_entrie", me)
}

type DataminingBlocklistEntry struct {
	Name string                  `json:"name"`
	Type DataminingBlocklistType `json:"type"` // Possible Values: `BUCKET`, `TABLE`
}

func (me *DataminingBlocklistEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BUCKET`, `TABLE`",
			Required:    true,
		},
	}
}

func (me *DataminingBlocklistEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name": me.Name,
		"type": me.Type,
	})
}

func (me *DataminingBlocklistEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name": &me.Name,
		"type": &me.Type,
	})
}
