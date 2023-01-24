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
	"sort"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResultMetadata struct {
	Entries []*ResultMetadataEntry
}

func (me *ResultMetadata) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"config": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Additional metadata for charted metric",
			Elem: &schema.Resource{
				Schema: new(ResultMetadataEntry).Schema(),
			},
		},
	}
}

func (me *ResultMetadata) MarshalHCL(properties hcl.Properties) error {
	if len(me.Entries) > 0 {
		entries := []any{}
		for _, entry := range me.Entries {
			marshalled := hcl.Properties{}
			if err := entry.MarshalHCL(marshalled); err == nil {
				entries = append(entries, marshalled)
			} else {
				return err
			}
		}
		sort.Slice(entries, func(i, j int) bool {
			d1, _ := json.Marshal(entries[i])
			d2, _ := json.Marshal(entries[j])
			cmp := strings.Compare(string(d1), string(d2))
			return (cmp == -1)
		})

		properties["config"] = entries
	}
	return nil
}

func (me *ResultMetadata) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("config.#"); ok {
		me.Entries = []*ResultMetadataEntry{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(ResultMetadataEntry)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "config", idx)); err != nil {
				return err
			}
			me.Entries = append(me.Entries, entry)
		}
	}
	return nil
}
