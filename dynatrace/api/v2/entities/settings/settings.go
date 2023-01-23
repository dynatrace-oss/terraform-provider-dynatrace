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

package entities

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	TotalCount  int64    `json:"totalCount"`            // The total number of entries in the result.
	PageSize    *int32   `json:"pageSize,omitempty"`    // The number of entries per page.
	NextPageKey *string  `json:"nextPageKey,omitempty"` // The cursor for the next page of results. Has the value of null on the last page.
	Entities    Entities `json:"entities,omitempty"`    // A list of monitored entities.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"total_count": {
			Type:        schema.TypeInt,
			Description: "The total number of entries in the result.",
			Required:    true,
		},
		"page_size": {
			Type:        schema.TypeInt,
			Description: "The number of entries per page.",
			Optional:    true,
		},
		"next_page_key": {
			Type:        schema.TypeString,
			Description: "The cursor for the next page of results. Has the value of null on the last page.",
			Optional:    true,
		},
		"entities": {
			Type:        schema.TypeList,
			Description: "A list of monitored entities.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Entities).Schema()},
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"total_count":   me.TotalCount,
		"page_size":     me.PageSize,
		"next_page_key": me.NextPageKey,
		"entities":      me.Entities,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"total_count":   &me.TotalCount,
		"page_size":     &me.PageSize,
		"next_page_key": &me.NextPageKey,
		"entities":      &me.Entities,
	})
}
