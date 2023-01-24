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

package common

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TagFilters []*TagFilter

func (me TagFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "A Tag Filter",
			Elem:        &schema.Resource{Schema: new(TagFilter).Schema()},
		},
	}
}

func (me TagFilters) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		entries := []any{}
		for _, entry := range me {
			marshalled := hcl.Properties{}
			if err := entry.MarshalHCL(marshalled); err == nil {
				entries = append(entries, marshalled)
			} else {
				return err
			}
		}
		if len(entries) > 0 {
			properties["filter"] = entries
		}
	}
	return nil
}

func (me *TagFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("filter"); ok {

		entrySet := value.(*schema.Set)

		for _, entryMap := range entrySet.List() {
			hash := entrySet.F(entryMap)
			entry := new(TagFilter)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "filter", hash)); err != nil {
				return err
			}
			*me = append(*me, entry)
		}
	}
	return nil
}
