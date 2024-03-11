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

package items

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HubItemList struct {
	Items       []*HubItem `json:"items"`
	NextPageKey *string    `json:"nextPageKey"`
	PageSize    *int32     `json:"pageSize"`
	TotalCount  *int64     `json:"totalCount"`
}

func (me *HubItemList) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"items": {
			Type:        schema.TypeList,
			Description: "The items within this list",
			Computed:    true,
			Elem:        &schema.Resource{Schema: new(HubItem).Schema()},
		},
	}
}

func (me *HubItemList) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("items", me.Items)
}

func (me *HubItemList) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("items", &me.Items)
}
